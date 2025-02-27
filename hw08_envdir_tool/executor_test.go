package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := Environment{
		"BAR":   EnvValue{Value: "bar"},
		"EMPTY": EnvValue{NeedRemove: true},
		"FOO":   EnvValue{Value: "   foo\nwith new line"},
		"HELLO": EnvValue{Value: "\"hello\""},
		"UNSET": EnvValue{NeedRemove: true},
	}

	t.Run("check define env vars", func(t *testing.T) {
		defineEnvVars(env)

		require.Equal(t, "bar", os.Getenv("BAR"))
		require.Equal(t, "", os.Getenv("EMPTY"))
		require.Equal(t, "   foo\nwith new line", os.Getenv("FOO"))
		require.Equal(t, "\"hello\"", os.Getenv("HELLO"))
		require.Equal(t, "", os.Getenv("UNSET"))
	})

	t.Run("check exec command", func(t *testing.T) {
		code := RunCmd([]string{"echo"}, env)
		require.Equal(t, 0, code)
	})
}
