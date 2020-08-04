package fileutils

import "strings"

func GetBaseFilename(filename string) string {
	lastExtension := strings.LastIndex(filename, ".")
	if lastExtension != -1 {
		return filename[:lastExtension]
	}
	return filename
}
