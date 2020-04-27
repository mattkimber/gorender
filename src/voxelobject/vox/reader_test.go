package vox

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
	"voxelobject"
)

func TestIsHeaderValid(t *testing.T) {
	var testCases = []struct {
		input    string
		expected bool
	}{
		{"VOX ", true},
		{"BLAH", false},
		{"A", false},
		{"ABCDE", false},
		{"VOX AAA", true},
	}

	for _, testCase := range testCases {
		reader := strings.NewReader(testCase.input)
		result := isHeaderValid(reader)
		if result != testCase.expected {
			t.Errorf("Magic string %s expected %v, got %v", testCase.input, testCase.expected, result)
		}

		expectedLength := len(testCase.input) - 4
		if expectedLength < 0 {
			expectedLength = 0
		}

		if reader.Len() != expectedLength {
			t.Errorf("Did not read 4 bytes of string %s, %d remaining of %d", testCase.input, reader.Len(), len(testCase.input))
		}
	}
}

func TestGetSizeFromChunk(t *testing.T) {
	var testCases = []struct {
		input    []byte
		expected voxelobject.Point
		error    bool
	}{
		{getSizedByteSlice(12, []byte{1, 0, 0, 0, 2, 0, 0, 0, 3, 0, 0, 0}),
			voxelobject.Point{X: 1, Y: 2, Z: 3}, false},
		{getSizedByteSlice(1, []byte{1}),
			voxelobject.Point{X: 0, Y: 0, Z: 0}, true},
		{getSizedByteSlice(200, []byte{1}),
			voxelobject.Point{X: 0, Y: 0, Z: 0}, true},
		{getSizedByteSlice(16, []byte{3, 0, 0, 0, 5, 0, 0, 0, 1, 0, 0, 0, 2, 4, 6, 8}),
			voxelobject.Point{X: 3, Y: 5, Z: 1}, false},
	}

	for _, testCase := range testCases {
		reader := bytes.NewReader(testCase.input)
		result, err := getSizeFromChunk(reader)

		if testCase.error && err == nil {
			t.Errorf("Expected error for input %v, got none", testCase.input)
		}

		if result != testCase.expected {
			t.Errorf("Byte array %v expected %v, got %v", testCase.input, testCase.expected, result)
		}

		if reader.Len() > 0 {
			t.Errorf("Did not read to end of data for byte array %v", testCase.input)
		}
	}
}

func TestGetPointDataFromChunk(t *testing.T) {
	var testCases = []struct {
		input    []byte
		expected []voxelobject.PointWithColour
		error bool
	}{
		{getSizedByteSlice(4, []byte{1, 2, 3, 64}),
			[]voxelobject.PointWithColour{{voxelobject.Point{X: 1, Y: 2, Z: 3}, 64}}, false},
		{getSizedByteSlice(8, []byte{1, 2, 3, 64, 4, 5, 6, 128}),
			[]voxelobject.PointWithColour{
				{voxelobject.Point{X: 1, Y: 2, Z: 3}, 64},
				{voxelobject.Point{X: 4, Y: 5, Z: 6}, 128},
			}, false},
		{getSizedByteSlice(5, []byte{1, 2, 3, 4, 5}),
			[]voxelobject.PointWithColour{}, true},

	}

	for _, testCase := range testCases {
		reader := bytes.NewReader(testCase.input)
		result, err := getPointDataFromChunk(reader)

		if testCase.error && err == nil {
			t.Errorf("Expected error for input %v, got none", testCase.input)
		}

		if !arePointWithColourSlicesEqual(result, testCase.expected) {
			t.Errorf("Byte array %v expected %v, got %v", testCase.input, testCase.expected, result)
		}

		if reader.Len() > 0 {
			t.Errorf("Did not read to end of data for byte array %v", testCase.input)
		}
	}
}

func TestSkipUnhandledChunk(t *testing.T) {
	var testCases = [][]byte {
		getSizedByteSlice(4, []byte{1,2,3,4}),
	}

	for _, testCase := range testCases {
		reader := bytes.NewReader(testCase)
		skipUnhandledChunk(reader)

		if reader.Len() > 0 {
			t.Errorf("Did not read to end of data for byte array %v", testCase)
		}
	}
}

func arePointWithColourSlicesEqual(a []voxelobject.PointWithColour, b []voxelobject.PointWithColour) bool {
	if len(a) != len(b) {
		return false
	}

	for i, p := range a {
		if p != b[i] {
			return false
		}
	}

	return true
}

func getSizedByteSlice(size int64, slice []byte) []byte {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, uint64(size))
	result = append(result, slice...)
	return result
}
