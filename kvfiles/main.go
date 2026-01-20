package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// get([]byte) []byte, error
// set(key, value []byte) error
// Есть интерфейс Get/Set. Нужно реализовать key-value хранилище на файлах.
// Многопоточный доступ. Можно использовать io и map.

type Storage struct {
	baseDir string
	inmem   map[string][]byte
	mx      sync.RWMutex
}

func newStorage(baseDir string) (*Storage, error) {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, err
	}
	return &Storage{
		baseDir: baseDir,
		inmem:   make(map[string][]byte),
	}, nil
}

func (s *Storage) filePath(key []byte) string {
	return filepath.Join(s.baseDir, string(key))
}

// -------------------- GET --------------------

func (s *Storage) get(key []byte) ([]byte, error) {
	k := string(key)

	// 1. Fast path — in-memory
	s.mx.RLock()
	val, ok := s.inmem[k]
	s.mx.RUnlock()
	if ok {
		cp := make([]byte, len(val))
		copy(cp, val)
		return cp, nil
	}

	// 2. Slow path — filesystem via io.Reader
	f, err := os.Open(s.filePath(key))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// 3. Cache in memory (double-check)
	s.mx.Lock()
	if cur, exists := s.inmem[k]; exists {
		cp := make([]byte, len(cur))
		copy(cp, cur)
		s.mx.Unlock()
		return cp, nil
	}
	cp := make([]byte, len(data))
	copy(cp, data)
	s.inmem[k] = cp
	s.mx.Unlock()

	return data, nil
}

// -------------------- SET --------------------

func (s *Storage) set(key, value []byte) error {
	// 1. Write to filesystem via io.Writer
	f, err := os.Create(s.filePath(key))
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(f)
	if _, err := writer.Write(value); err != nil {
		_ = f.Close()
		return err
	}
	if err := writer.Flush(); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	// 2. Update in-memory cache
	cp := make([]byte, len(value))
	copy(cp, value)
	s.mx.Lock()
	s.inmem[string(key)] = cp
	s.mx.Unlock()
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	baseDir := "./kvdata"
	_ = os.RemoveAll(baseDir)

	s, err := newStorage(baseDir)
	if err != nil {
		panic(err)
	}

	// Тестовые параметры
	const (
		keysCount    = 16
		writers      = 8
		readers      = 16
		opsPerWorker = 5000
	)

	// Глобальная "правда" для проверки
	var expected sync.Map // key string -> []byte

	// Инициализация ключей
	keys := make([]string, 0, keysCount)
	for i := 0; i < keysCount; i++ {
		k := fmt.Sprintf("k%02d.txt", i)
		keys = append(keys, k)
		v := []byte("init_" + k)
		if err := s.set([]byte(k), v); err != nil {
			panic(err)
		}
		expected.Store(k, append([]byte(nil), v...))
	}

	var writeOps atomic.Int64
	var readOps atomic.Int64
	var mismatches atomic.Int64

	// Барьер старта
	start := make(chan struct{})
	var wg sync.WaitGroup

	// Writers
	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-start
			for i := 0; i < opsPerWorker; i++ {
				k := keys[rand.Intn(len(keys))]
				v := []byte("w" + strconv.Itoa(id) + "_i" + strconv.Itoa(i))
				if err := s.set([]byte(k), v); err != nil {
					panic(err)
				}
				expected.Store(k, append([]byte(nil), v...))
				writeOps.Add(1)
				if i%200 == 0 {
					time.Sleep(time.Millisecond)
				}
			}
		}(w)
	}

	// Readers
	for r := 0; r < readers; r++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-start
			for i := 0; i < opsPerWorker; i++ {
				k := keys[rand.Intn(len(keys))]
				v, err := s.get([]byte(k))
				if err != nil {
					panic(err)
				}
				if v == nil {
					mismatches.Add(1)
				}
				readOps.Add(1)
				if i%500 == 0 {
					time.Sleep(time.Millisecond)
				}
			}
		}(r)
	}

	close(start)
	wg.Wait()

	// Фаза стабилизации
	for _, k := range keys {
		wantAny, _ := expected.Load(k)
		want := wantAny.([]byte)
		got, err := s.get([]byte(k))
		if err != nil {
			panic(err)
		}
		if !bytes.Equal(got, want) {
			fmt.Printf("FINAL MISMATCH key=%s got=%q want=%q\n", k, string(got), string(want))
			mismatches.Add(1)
		}
	}

	fmt.Println("DONE")
	fmt.Println("writeOps:", writeOps.Load())
	fmt.Println("readOps :", readOps.Load())
	fmt.Println("mismatches:", mismatches.Load())
	fmt.Printf("inmem keys: %d\n", func() int {
		s.mx.Lock()
		defer s.mx.Unlock()
		return len(s.inmem)
	}())
}
