package filesystem

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// unzipToDir https://stackoverflow.com/a/24792688/11119471
func unzipToDir(zipContent io.ReaderAt, assumedSize int64, destPath string) (bytesWritten int, err error) {
	zipR, err := zip.NewReader(zipContent, assumedSize)
	if err != nil {
		return 0, err
	}
	if err = os.MkdirAll(destPath, 0755); err != nil {
		return 0, err
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) (n int, err error) {
		rc, err := f.Open()
		if err != nil {
			return 0, err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(destPath, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(destPath)+string(os.PathSeparator)) {
			return 0, fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return 0, err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			nn, err := io.Copy(f, rc)
			n = int(nn)
			if err != nil {
				return 0, err
			}
		}
		return n, nil
	}

	for _, f := range zipR.File {
		n, err := extractAndWriteFile(f)
		if err != nil {
			return bytesWritten, err
		}
		bytesWritten += n
	}

	return bytesWritten, nil
}

// zipDir https://stackoverflow.com/a/49057861/11119471
func zipDir(path string) ([]byte, error) {
	destination := &bytes.Buffer{}
	myZip := zip.NewWriter(destination)

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if info == nil {
			return fmt.Errorf("fileinfo [%s] is nil, %w", filePath, os.ErrNotExist)
		}
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, path)
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	err = myZip.Close()
	if err != nil {
		return nil, err
	}
	return destination.Bytes(), nil
}

func removeDirContent(path string) error {
	glob := filepath.Join(path, "*")
	files, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	for _, f := range files {
		if err = os.RemoveAll(f); err != nil {
			return err
		}
	}
	return nil
}
