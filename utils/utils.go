package utils

import "path/filepath"

func GetExtensionFile(fileName string) string {
	return filepath.Ext(fileName)
}
