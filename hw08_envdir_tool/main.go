package main

import "os"

func main() {
	paramsCount := len(os.Args)

	if paramsCount == 1 {
		panic("Directory params is not defined")
	}

	if paramsCount == 2 {
		panic("Command params is not defined")
	}

	directory := os.Args[1]
	commandParams := os.Args[2:]

	envVariables, err := ReadDir(directory)
	if err != nil {
		panic(err)
	}

	RunCmd(commandParams, envVariables)
}
