package logger

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	filePath := "log"

	removeLogFile := func() error {
		return os.Remove(filePath)
	}

	t.Run("correct level", func(t *testing.T) {
		defer removeLogFile()

		levels := []string{"error", "warn", "info", "debug", "ERROR", "WARN", "INFO", "DEBUG"}
		for _, level := range levels {
			_, err := New(level, filePath)
			require.NoError(t, err)
		}
	})

	t.Run("incorrect level", func(t *testing.T) {
		_, err := New("incorrect_level", filePath)
		require.Error(t, err)
	})

	t.Run("create log file", func(t *testing.T) {
		defer removeLogFile()

		logger, err := New("debug", filePath)
		if err != nil {
			log.Fatalln(err)
		}

		logger.Info("info message")
		logger.Debug("debug message")
		logger.Warning("warning message")
		logger.Error("error message")

		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		require.NoError(t, err)

		content := string(data)
		lines := strings.Split(content, `\n`)

		require.Equal(t, 4, len(lines))
	})
}
