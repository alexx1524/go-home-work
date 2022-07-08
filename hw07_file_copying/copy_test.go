package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	source := "testdata/input.txt"
	destination := "dest.txt"
	defer os.Remove(destination)

	tests := []struct {
		text      string
		offset    int64
		limit     int64
		source    string
		checkFile string
		err       error
	}{
		{offset: 0, limit: 0, source: source, checkFile: "testdata/out_offset0_limit0.txt"},
		{offset: 0, limit: 10, source: source, checkFile: "testdata/out_offset0_limit10.txt"},
		{offset: 0, limit: 1000, source: source, checkFile: "testdata/out_offset0_limit1000.txt"},
		{offset: 100, limit: 1000, source: source, checkFile: "testdata/out_offset100_limit1000.txt"},
		{offset: 0, limit: 10000, source: source, checkFile: "testdata/out_offset0_limit10000.txt"},
		{offset: 6000, limit: 1000, source: source, checkFile: "testdata/out_offset6000_limit1000.txt"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("Offset: %v, limit %v %s", tc.offset, tc.limit, tc.text), func(t *testing.T) {
			err := Copy(source, destination, tc.offset, tc.limit)

			if tc.err == nil {
				require.Nil(t, err)

				srcContent, _ := ioutil.ReadFile(tc.checkFile)
				destContent, _ := ioutil.ReadFile(destination)

				require.Equal(t, string(srcContent), string(destContent))
			} else {
				require.Error(t, tc.err, err)
			}
		})
	}

	t.Run("Offset exceeds file size", func(t *testing.T) {
		err := Copy(source, destination, 10000, 0)
		require.Error(t, err, ErrFromPathIsUndefined)
	})

	t.Run("FromPath is undefined", func(t *testing.T) {
		err := Copy("", destination, 0, 0)
		require.Error(t, err, ErrFromPathIsUndefined)
	})

	t.Run("ToPath is undefined", func(t *testing.T) {
		err := Copy(source, "", 0, 0)
		require.Error(t, err, ErrToPathIsUndefined)
	})

	t.Run("Offset is negative", func(t *testing.T) {
		err := Copy(source, destination, -1, 0)
		require.Error(t, err, ErrOffsetIsNegative)
	})

	t.Run("Limit is negative", func(t *testing.T) {
		err := Copy(source, destination, 0, -1)
		require.Error(t, err, ErrOffsetIsNegative)
	})

	t.Run("Unsupported source file", func(t *testing.T) {
		err := Copy("/dev/urandom", destination, 0, 0)

		require.Equal(t, ErrUnsupportedFile, err)
	})
}
