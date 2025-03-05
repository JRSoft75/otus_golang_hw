package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Test with environment variables", func(t *testing.T) {
		dir := "./testdata/env"
		expectedEnv := Environment{
			"BAR": EnvValue{
				"bar",
				false,
			},
			"EMPTY": EnvValue{
				"",
				false,
			},
			"FOO": EnvValue{
				"   foo\nwith new line",
				false,
			},
			"HELLO": EnvValue{
				"\"hello\"",
				false,
			},
			"UNSET": EnvValue{
				"",
				true,
			},
		}

		result, err := ReadDir(dir)
		require.Equal(t, err, nil)
		require.Equal(t, expectedEnv, result)
	})
}
