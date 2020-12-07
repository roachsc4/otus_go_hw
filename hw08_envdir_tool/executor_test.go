package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("success echo", func(t *testing.T) {
		cmd := []string{"testdata/write_file.sh", "arg1=1", "arg2=2"}
		environment := Environment{
			"BAR":   "bar",
			"FOO":   "   foo\nwith new line",
			"HELLO": `"hello"`,
			"UNSET": "",
		}
		returnCode := RunCmd(cmd, environment)
		require.Equal(t, 0, returnCode)
	})

	t.Run("wrong command string", func(t *testing.T) {
		for _, cmd := range [][]string{
			{},
			{"command"},
		} {
			environment := Environment{}
			returnCode := RunCmd(cmd, environment)
			require.Equal(t, InvalidCommandString, returnCode)
		}
	})

	t.Run("using env with empty key", func(t *testing.T) {
		environment := Environment{
			"": "value",
		}
		returnCode := RunCmd([]string{"command", "arg"}, environment)
		require.Equal(t, -1, returnCode)
	})
}
