package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Test with environment variables", func(t *testing.T) {
		cmd := []string{
			"testdata/echo.sh",
			"arg1=1",
			"arg2=2",
		}

		testEnv := Environment{
			"BAR": EnvValue{
				"bar",
				false,
			},
			"FOO": EnvValue{
				"   foo\nwith new line",
				false,
			},
			"EMPTY": EnvValue{
				"",
				false,
			},
		}

		result := RunCmd(cmd, testEnv)
		require.Equal(t, 0, result)
	})
}
