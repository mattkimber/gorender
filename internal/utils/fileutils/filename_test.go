package fileutils

import "testing"

func TestGetBaseFilename(t *testing.T) {
	testCases := []struct {
		input, expected string
	}{
		{"test", "test"},
		{"test.png", "test"},
		{"files/test.png", "files/test"},
		{"test.a.b.c", "test.a.b"},
		{".render/a/b/c", ".render/a/b/c"},
		{".render/a/b/c.png", ".render/a/b/c"},
	}

	for _, testCase := range testCases {
		if result := GetBaseFilename(testCase.input); result != testCase.expected {
			t.Errorf("input %s got %s, expected %s", testCase.input, result, testCase.expected)
		}
	}
}
