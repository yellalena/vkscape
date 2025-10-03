package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	OutputDir      = "vkscape_output"
	OutputGroupDir = "group_%s"
)

func CreateGroupDirectory(groupID string) string {
	groupDir := filepath.Join(OutputDir, fmt.Sprintf(OutputGroupDir, groupID))
	_ = os.MkdirAll(groupDir, 0755)
	return groupDir
}

func CreateSubDirectory(parentDir, subDir string) string {
	dir := filepath.Join(parentDir, subDir)
	_ = os.MkdirAll(dir, 0755)
	return dir
}

func SaveFile(parentDir, filename string, content []byte) error {
	filePath := filepath.Join(parentDir, filename)
	return os.WriteFile(filePath, content, 0644)
}

func SaveObject(parentDir, filename string, content io.ReadCloser) error {
	filePath := filepath.Join(parentDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, content)
	return err
}
