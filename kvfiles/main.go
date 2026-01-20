package main

import (
	"errors"
	"fmt"
	"os"
)

// get([]byte) []byte, error
// set(key, value []byte) error
// начнем короче с
// Есть интерфейс Get/Set. Нужно реализовать key-value хранилище на файлах.
// Многопоточный доступ. Можно использовать io и map

type Storage struct {
	inmem map[string][]byte
}

func newStorage() Storage {
	// add initialization from file system
	return Storage{
		inmem: make(map[string][]byte),
	}
}

func (s Storage) get(key []byte) ([]byte, error) {
	val, ok := s.inmem[string(key)]
	if ok {
		return val, nil
	}

	var err error
	val, err = os.ReadFile(string(key))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	s.inmem[string(key)] = val

	return val, nil
}

func (s Storage) set(key, value []byte) error {
	err := os.WriteFile(string(key), value, 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	s := newStorage()

	data, _ := s.get([]byte("a.txt"))
	fmt.Println(string(data))

	_ = s.set([]byte("a.txt"), []byte("bbb"))

	data, _ = s.get([]byte("a.txt"))
	fmt.Println(string(data))

	fmt.Printf("%+v", s)
}
