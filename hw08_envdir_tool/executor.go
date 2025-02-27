package main

import (
	"fmt"
	"os"
	"os/exec"
)

func defineEnvVars(env Environment) {
	for name, v := range env {
		if err := os.Unsetenv(name); err != nil {
			fmt.Printf("couldn't unset environment variable %s: %v\n", name, err)
		}
		if !v.NeedRemove {
			if err := os.Setenv(name, v.Value); err != nil {
				fmt.Printf("couldn't set environment variable %s: %v\n", name, err)
			}
		}
	}
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		fmt.Println("wrong command")
		return 1
	}

	cmdName := cmd[0]
	cmdArgs := cmd[1:]

	defineEnvVars(env)

	command := exec.Command(cmdName, cmdArgs...)
	command.Env = os.Environ()
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	if err := command.Run(); err != nil {
		return 1
	}

	return 0
}
