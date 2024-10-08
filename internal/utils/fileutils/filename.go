package fileutils

import "strings"

func GetBaseFilename(filename string) string {
	lastExtension := strings.LastIndex(filename, ".")
	lastSlash := strings.LastIndex(filename, "/")
	if lastExtension != -1 && lastExtension > lastSlash {
		return filename[:lastExtension]
	}
	return filename
}
