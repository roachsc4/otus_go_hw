package main

import (
	"errors"
	"os"
	"os/exec"
)

const InvalidCommandString = 418

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 2 {
		return InvalidCommandString
	}
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	for variable, value := range env {
		if value == "" {
			err := os.Unsetenv(variable)
			if err != nil {
				return -1
			}
		}
		err := os.Setenv(variable, value)
		if err != nil {
			return -1
		}
	}

	command.Env = os.Environ()
	err := command.Run()

	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return exitError.ExitCode()
	}
	return 0
}
