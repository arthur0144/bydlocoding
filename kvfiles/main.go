package main

import (
	"fmt"
	"os"
)

// get([]byte) []byte, error
// set(key, value []byte) error
// начнем короче с
// Есть интерфейс Get/Set. Нужно реализовать key-value хранилище на файлах.
// Многопоточный доступ. Можно использовать io и map

type Storage struct {
}

// add init method

func (s Storage) get(key []byte) ([]byte, error) {
	val, err := os.ReadFile(string(key))
	if err != nil {
		return nil, err
	}
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
	var s Storage

	data, _ := s.get([]byte("a.txt"))
	fmt.Println(string(data))

	_ = s.set([]byte("a.txt"), []byte("bbb"))

	data, _ = s.get([]byte("a.txt"))
	fmt.Println(string(data))

	fmt.Printf("%+v", s)
}
