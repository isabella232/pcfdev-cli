package fs

import (
	"fmt"
	"io"
	"os"
)

type FS struct{}

func (fs *FS) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (fs *FS) Write(path string, contents io.ReadCloser) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}
	defer file.Close()
	_, err = io.Copy(file, contents)
	if err != nil {
		return fmt.Errorf("failed to copy contents to file: %s", err)
	}
	return nil
}

func (fs *FS) CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}
