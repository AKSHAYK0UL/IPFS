package fileop

import (
	"os"
)

type FilePointer struct {
	File *os.File `json:"file"`
}

// custom write method
func (fp *FilePointer) write(data []byte) error {
	_, err := fp.File.Write(data)
	return err
}
