package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("check result env map", func(t *testing.T) {
		dir := "testdata/env"

		env, err := ReadDir(dir)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		expectedEnv := Environment{
			"BAR":   EnvValue{Value: "bar"},
			"EMPTY": EnvValue{NeedRemove: true},
			"FOO":   EnvValue{Value: "   foo\nwith new line"},
			"HELLO": EnvValue{Value: "\"hello\""},
			"UNSET": EnvValue{NeedRemove: true},
		}
		require.Equal(t, expectedEnv, env)
	})
}
