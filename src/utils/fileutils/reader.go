package fileutils

import (
	"io"
	"os"
)

type FileReader interface {
	GetFromReader(r io.Reader) error
}

type reader struct {
	fileReader FileReader
}

func (r reader) GetFileHandle(filename string) (f *os.File, err error) {
	f, err = os.Open(filename)
	return
}

func (r reader) DoIO(f *os.File) (err error) {
	err = r.fileReader.GetFromReader(f)
	return
}
