package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := make(Environment)

	for _, file := range files {
		fileName := file.Name()
		value, err := getValue(path.Join(dir, fileName))
		if err != nil {
			return nil, err
		}

		result[fileName] = value
	}

	return result, err
}

func getValue(fileName string) (EnvValue, error) {
	result := EnvValue{}

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return result, err
	}
	size := fileInfo.Size()
	if size == 0 {
		result.NeedRemove = true
		return result, nil
	}

	readFile, err := os.Open(fileName)
	defer readFile.Close()

	if err != nil {
		return result, err
	}
	fileScanner := bufio.NewScanner(readFile)
	if fileScanner.Scan() {
		line := fileScanner.Text()
		line = strings.TrimRight(line, string(rune(0x00)))
		line = strings.TrimRight(line, "\t")
		line = strings.TrimRight(line, " ")
		line = strings.ReplaceAll(line, "\x00", "\n")

		if line == "" {
			result.NeedRemove = true
			return result, nil
		}
		result.Value = line
	}

	return result, nil
}
