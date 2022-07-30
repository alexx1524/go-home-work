package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	getEnviron := func(env Environment) ([]string, error) {
		if env == nil {
			return os.Environ(), nil
		}
		for name, value := range env {
			if value.NeedRemove {
				err := os.Unsetenv(name)
				if err != nil {
					return nil, err
				}
			} else {
				err := os.Setenv(name, value.Value)
				if err != nil {
					return nil, err
				}
			}
		}
		return os.Environ(), nil
	}

	commandName := cmd[0]
	command := exec.Command(commandName, cmd[1:]...)

	environments, err := getEnviron(env)
	if err != nil {
		println("Getting environment variables error: ", err)
	}

	command.Env = environments
	command.Stdout = os.Stdout

	if err := command.Run(); err != nil {
		println("Executing error: ", err)
	}

	return command.ProcessState.ExitCode()
}
