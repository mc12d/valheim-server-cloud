package filesystem_test

import (
	"os"
	"regexp"
	"testing"

	"backup-agent/internal/adapters/filesystem"

	"github.com/stretchr/testify/require"
)

func TestZipUnzip(t *testing.T) {
	defer os.RemoveAll("testdata/testdir_result")
	zipbytes, err := filesystem.DirectoryRegexZipper(regexp.MustCompile(".*\\.txt$"))("testdata/testdir")
	require.NoError(t, err)

	err = filesystem.DirectoryUnzipper("testdata/testdir_result", zipbytes)
	require.NoError(t, err)

	require.DirExists(t, "testdata/testdir_result")
	require.FileExists(t, "testdata/testdir_result/hello.txt")
	require.FileExists(t, "testdata/testdir_result/inner/world.txt")
	require.NoFileExists(t, "testdata/testdir_result/notxt.sh")
	require.NoFileExists(t, "testdata/testdir_result/wld.txt.old")

}
