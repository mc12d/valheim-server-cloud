package filesystem

import (
	"backup-agent/internal/app"
	"bytes"
	"regexp"
)

var (
	DirectoryZipper app.BackupZipper = func(path string) ([]byte, error) {
		return zipDir(path)
	}

	DirectoryUnzipper app.BackupUnzipper = func(path string, zip []byte) error {
		if err := removeDirContent(path); err != nil {
			return err
		}
		_, err := unzipToDir(bytes.NewReader(zip), int64(len(zip)), path)
		return err
	}
)

func DirectoryRegexZipper(dirContentRegex *regexp.Regexp) app.BackupZipper {
	return func(path string) ([]byte, error) {
		return zipDirRegex(path, dirContentRegex)
	}
}
