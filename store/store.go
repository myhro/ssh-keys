package store

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

const perm = 0o600

type File struct {
	Path string
}

func New(path string) *File {
	return &File{Path: path}
}

func (f *File) Read() ([]byte, error) {
	data, err := os.ReadFile(f.Path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", f.Path, err)
	}
	return data, nil
}

func (f *File) Write(data []byte) error {
	tmp := f.Path + ".tmp"
	defer os.Remove(tmp)

	err := os.WriteFile(tmp, data, perm)
	if err != nil {
		return fmt.Errorf("writing %s: %w", tmp, err)
	}

	err = os.Rename(tmp, f.Path)
	if err != nil {
		return fmt.Errorf("renaming %s: %w", tmp, err)
	}

	return nil
}
