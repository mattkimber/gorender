package vox

import (
	"strings"
	"testing"
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
