package fileutils

import (
	"os"
)

type fileIOHandler interface {
	GetFileHandle(filename string) (*os.File, error)
	DoIO(f *os.File) error
}

func InstantiateFromFile(filename string, o FileReader) (err error) {
	r := reader{o}
	err = doFileIO(filename, r)
	return
}

func WriteToFile(filename string, o FileWriter) (err error) {
	w := writer{o}
	err = doFileIO(filename, w)
	return
}

func doFileIO(filename string, handler fileIOHandler) (err error) {
	file, err := handler.GetFileHandle(filename)
	if err != nil {
		return
	}

	err = handler.DoIO(file)

	if err != nil {
		file.Close()
		return
	}

	if err = file.Close(); err != nil {
		return
	}

	return
}
