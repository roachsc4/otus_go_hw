package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		require.NoError(t, err)
		expectedEnv := Environment{
			"BAR":   "bar",
			"FOO":   "   foo\nwith new line",
			"HELLO": `"hello"`,
			"UNSET": "",
		}
		require.Equal(t, expectedEnv, env)
	})

	t.Run("non-existing directory", func(t *testing.T) {
		_, err := ReadDir("foo")
		require.EqualError(t, err, "failed to read directory data: open foo: no such file or directory")
	})

	t.Run("open dir with wrong env-file", func(t *testing.T) {
		_, err := ReadDir("testdata/env_with_wrong_file")
		require.EqualError(t, err, "failed to process WRONG=FILE: it's name contains '='")
	})
}
