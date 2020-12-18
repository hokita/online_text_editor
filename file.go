package main

import (
	"io/ioutil"
)

type file struct {
	name string
}

func newFile(filename string) *file {
	return &file{name: filename}
}

func (f *file) read() (string, error) {
	byte, err := ioutil.ReadFile(f.name)
	if err != nil {
		return "", err
	}

	return string(byte), nil
}

func (f *file) write(b []byte) error {
	err := ioutil.WriteFile(f.name, b, 0666)
	if err != nil {
		return err
	}
	return nil
}
