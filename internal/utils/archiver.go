//nolint:revive // we're fine with utils
package utils

// import (
// 	"archive/zip"
// 	"io"
// 	"os"
// 	"path/filepath"
// 	"strings"
// )

// func ZipFolder(filename string) error {
// 	// Archive whole output directory
// 	sourceDir := OutputDir

// 	archivePath := filepath.Join(sourceDir, filename)
// 	zipFile, err := os.Create(archivePath)
// 	if err != nil {
// 		return err
// 	}
// 	defer zipFile.Close()

// 	writer := zip.NewWriter(zipFile)
// 	defer writer.Close()

// 	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		if path == archivePath {
// 			return nil
// 		}

// 		// Only keep relative path — assume it's already normalized and clean
// 		relPath := strings.TrimPrefix(path, sourceDir)
// 		relPath = strings.TrimLeft(relPath, string(filepath.Separator)) // ensure no leading slash

// 		if relPath == "" {
// 			return nil
// 		}

// 		header, err := zip.FileInfoHeader(info)
// 		if err != nil {
// 			return err
// 		}
// 		header.Name = filepath.ToSlash(relPath)

// 		if info.IsDir() {
// 			header.Name += "/"
// 			_, err := writer.CreateHeader(header)
// 			return err
// 		}

// 		header.Method = zip.Deflate
// 		writer, err := writer.CreateHeader(header)
// 		if err != nil {
// 			return err
// 		}

// 		file, err := os.Open(path)
// 		if err != nil {
// 			return err
// 		}
// 		defer file.Close()

// 		_, err = io.Copy(writer, file)
// 		return err
// 	})
// }
