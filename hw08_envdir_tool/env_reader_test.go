package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	t.Run("Success read", func(t *testing.T) {
		environment, err := ReadDir("testdata/env")
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, len(environment), 5)

		require.Equal(t, environment["BAR"].Value, "bar")
		require.False(t, environment["BAR"].NeedRemove)

		require.Empty(t, environment["EMPTY"].Value)
		require.True(t, environment["EMPTY"].NeedRemove)

		require.Equal(t, environment["FOO"].Value, "   foo\nwith new line")
		require.False(t, environment["FOO"].NeedRemove)

		require.Equal(t, environment["HELLO"].Value, `"hello"`)
		require.False(t, environment["HELLO"].NeedRemove)

		require.Empty(t, environment["UNSET"].Value)
		require.True(t, environment["UNSET"].NeedRemove)
	})

	t.Run("If directory doesn't exist the function returns error", func(t *testing.T) {
		_, err := ReadDir("wrong_directory")
		require.Error(t, err)
	})

	t.Run("If directory is empty the function returns empty map", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "temp")
		if err != nil {
			log.Fatalln(err)
		}
		defer os.RemoveAll(dir)

		environment, err := ReadDir(dir)
		require.Equal(t, len(environment), 0)
	})

	t.Run("if file is empty the function returns need remove the environment variable", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "test")
		if err != nil {
			log.Fatalln(err)
		}

		file, err := ioutil.TempFile(dir, "test.txt")
		if err != nil {
			log.Fatalln(err)
		}

		defer os.Remove(file.Name())
		defer os.RemoveAll(dir)

		environment, err := ReadDir(dir)

		require.Equal(t, len(environment), 1)
		require.True(t, environment[filepath.Base(file.Name())].NeedRemove)
	})
}
