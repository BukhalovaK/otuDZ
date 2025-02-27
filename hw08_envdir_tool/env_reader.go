package main

import (
	"bufio"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return env, err
	}

	for _, entry := range entries {
		file, err := os.Open(dir + "/" + entry.Name())
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		firstLine := ""
		if scanner.Scan() {
			firstLine = strings.ReplaceAll(scanner.Text(), "\x00", "\n")
			firstLine = strings.TrimRight(firstLine, " ")
		} else if err := scanner.Err(); err != nil {
			return nil, err
		}

		envVal := EnvValue{}
		envVal.Value = firstLine
		if envVal.Value == "" {
			envVal.NeedRemove = true
		}
		env[entry.Name()] = envVal
	}
	return env, nil
}
