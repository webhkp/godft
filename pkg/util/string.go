package util

import (
	"path"
	"path/filepath"
)

func GetFileNameWithoutExtension(fullPath string) string {
	filenameWithoutExt := path.Base(fullPath)

	ext := filepath.Ext(fullPath)

	filename := filenameWithoutExt[:len(filenameWithoutExt)-len(ext)]

	return filename
}
