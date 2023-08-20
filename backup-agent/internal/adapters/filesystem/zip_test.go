package filesystem_test

import (
	"testing"

	"backup-agent/internal/adapters/filesystem"

	"github.com/stretchr/testify/require"
)

func TestZipUnzip(t *testing.T) {
	// defer os.RemoveAll("testdata/testdir_result")
	zipbytes, err := filesystem.ZipDirectory("testdata/testdir")
	require.NoError(t, err)
	require.True(t, len(zipbytes) > 1000, "actual len: {}", len(zipbytes))

	err = filesystem.UnzipToDirectory("testdata/testdir_result", zipbytes)
	require.NoError(t, err)

	require.DirExists(t, "testdata/testdir_result")
	require.FileExists(t, "testdata/testdir_result/hello.txt")
	require.FileExists(t, "testdata/testdir_result/inner/world.txt")

}
