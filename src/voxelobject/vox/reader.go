package vox

import (
	"io"
	"io/ioutil"
)

const magic = "VOX "

func isHeaderValid(handle io.Reader) bool {
	limitedReader := io.LimitReader(handle, 4)
	result, err := ioutil.ReadAll(limitedReader)
	return err == nil && string(result) == magic
}
