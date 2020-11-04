package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	limit         int64
	offset        int64
	checkFileName string
}

func TestCopy(t *testing.T) {
	for _, tst := range []test{
		{
			limit:         0,
			offset:        0,
			checkFileName: "testdata/out_offset0_limit0.txt",
		},
		{
			limit:         10,
			offset:        0,
			checkFileName: "testdata/out_offset0_limit10.txt",
		},
		{
			limit:         1000,
			offset:        0,
			checkFileName: "testdata/out_offset0_limit1000.txt",
		},
		{
			limit:         10000,
			offset:        0,
			checkFileName: "testdata/out_offset0_limit10000.txt",
		},
		{
			limit:         1000,
			offset:        100,
			checkFileName: "testdata/out_offset100_limit1000.txt",
		},
		{
			limit:         1000,
			offset:        6000,
			checkFileName: "testdata/out_offset6000_limit1000.txt",
		},
	} {
		t.Run(fmt.Sprintf("Limit: %d, Offset: %d", tst.limit, tst.offset), func(t *testing.T) {
			err := Copy("testdata/input.txt", "/tmp/output.txt", tst.offset, tst.limit)
			require.NoError(t, err)

			targetFile, err := os.Open("/tmp/output.txt")
			defer os.Remove(targetFile.Name())
			defer targetFile.Close()
			require.NoError(t, err)

			checkFile, _ := os.Open(tst.checkFileName)
			defer checkFile.Close()

			checkFileContent, _ := ioutil.ReadAll(checkFile)
			targetFileContent, _ := ioutil.ReadAll(targetFile)
			require.Equal(t, checkFileContent, targetFileContent)
		})
	}

	t.Run("offset bigger than file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/tmp/output.txt", 999999, 0)
		require.Equal(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("copy from non-existing file", func(t *testing.T) {
		err := Copy("no_file", "/tmp/output.txt", 999999, 0)
		require.Error(t, err)
	})

	t.Run("copy from directory", func(t *testing.T) {
		err := Copy("testdata", "/tmp/output.txt", 999999, 0)
		require.Equal(t, err, ErrUnsupportedFile)
	})

	t.Run("copy to file in non-existing directory", func(t *testing.T) {
		err := Copy("testdata/input.txt", "some/output.txt", 0, 0)
		require.Error(t, err)
	})
}
