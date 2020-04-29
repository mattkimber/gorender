package fileutils

import (
	"io"
	"os"
)

type FileWriter interface {
	OutputToWriter(w io.Writer) error
}

type writer struct {
	fileWriter FileWriter
}

func (w writer) GetFileHandle(filename string) (f *os.File, err error) {
	f, err = os.Create(filename)
	return
}

func (w writer) DoIO(f *os.File) (err error) {
	err = w.fileWriter.OutputToWriter(f)
	return
}
