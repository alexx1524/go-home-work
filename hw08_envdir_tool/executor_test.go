package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("If command is not found then RunCmd returns error code", func(t *testing.T) {
		cmd := []string{
			"/bin/bash",
			"testdata/echo2.sh",
		}

		returnCode := RunCmd(cmd, nil)

		require.Equal(t, returnCode, 127)
	})

	t.Run("If running data is correct then RunCmd returns 0 code", func(t *testing.T) {
		data := []string{
			"/bin/bash",
			"testdata/echo.sh",
		}

		returnCode := RunCmd(data, nil)

		require.Equal(t, returnCode, 0)
	})
}
