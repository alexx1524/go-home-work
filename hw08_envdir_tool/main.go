package main

import (
	"log"
	"os"
)

func main() {
	paramsCount := len(os.Args)

	if paramsCount == 1 {
		log.Fatalln("Directory params is not defined")
	}

	if paramsCount == 2 {
		log.Fatalln("Command params is not defined")
	}

	directory := os.Args[1]
	commandParams := os.Args[2:]

	envVariables, err := ReadDir(directory)
	if err != nil {
		log.Fatalln("Ошибка чтения переменных окружения из директории", err)
	}

	returnCode := RunCmd(commandParams, envVariables)

	os.Exit(returnCode)
}
