package fileutils

import (
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

type testData struct {
	value string
}

func (d *testData) GetFromReader(r io.Reader) error {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	d.value = string(bytes)
	return nil
}

func (d *testData) OutputToWriter(w io.Writer) (err error) {
	_, err = w.Write([]byte(d.value))
	return
}

func TestInstantiateFromFile(t *testing.T) {
	var f testData

	if err := InstantiateFromFile("testdata/sample.txt", &f); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if f.value != "hello world" {
		t.Errorf("expected 'hello world' but was '%s'", f.value)
	}
}

func TestInstantiateFromFile_FileNotExists(t *testing.T) {
	var f testData

	if err := InstantiateFromFile("testdata/not_exists.txt", &f); err == nil {
		t.Errorf("expected error, found none")
	}
}

func TestWriteToFile(t *testing.T) {
	const expected = "hello world"
	td := testData{expected}
	filename := "testdata/" + strconv.Itoa(int(rand.Uint64())) + ".txt"

	if err := WriteToFile(filename, &td); err != nil {
		t.Errorf("unexpected error writing: %v", err)
	}

	var f testData

	if err := InstantiateFromFile(filename, &f); err != nil {
		t.Errorf("unexpected error reading: %v", err)
	}

	if f.value != expected {
		t.Errorf("expected '%s', got '%s'", expected, f.value)
	}

	if err := os.Remove(filename); err != nil {
		t.Errorf("unexpected error deleting: %v", err)
	}
}
