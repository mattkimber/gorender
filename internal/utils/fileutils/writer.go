package fileutils

import (
	"bufio"
	"io"
	"os"
)

const writerSize = 1024 * 32

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
	buf := bufio.NewWriterSize(f, writerSize)
	err = w.fileWriter.OutputToWriter(buf)
	_ = buf.Flush()
	return
}
