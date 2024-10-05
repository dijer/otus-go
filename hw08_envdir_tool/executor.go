package main

import (
	"os"
	"os/exec"
)

const (
	errCode     = 1
	successCode = 0
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for envName, envValue := range env {
		if envValue.NeedRemove {
			err := os.Unsetenv(envName)
			if err != nil {
				return errCode
			}
			continue
		}

		err := os.Setenv(envName, envValue.Value)
		if err != nil {
			return errCode
		}
	}

	cmdCommand, cmdArgs := cmd[0], cmd[1:]
	command := exec.Command(cmdCommand, cmdArgs...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	err := command.Run()
	if err != nil {
		return command.ProcessState.ExitCode()
	}

	return successCode
}
