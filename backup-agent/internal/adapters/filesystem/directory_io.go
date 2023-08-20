package filesystem

import (
	"bytes"
)

type DirWriteMode int

const (
	Mb int64 = 2 << 20
)

func ZipDirectory(path string) (zip []byte, err error) {
	return zipDir(path)
}

func UnzipToDirectory(path string, zip []byte) error {
	if err := removeDirContent(path); err != nil {
		return err
	}
	_, err := unzipToDir(bytes.NewReader(zip), int64(len(zip)), path)
	return err
}
