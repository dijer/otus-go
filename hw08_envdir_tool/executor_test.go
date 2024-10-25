package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("should rewrite env", func(t *testing.T) {
		os.Setenv("BAR", "foo")

		returnCode := RunCmd([]string{"echo", ""}, Environment{
			"BAR": EnvValue{"bar", false},
		})

		require.Equal(t, os.Getenv("BAR"), "bar")
		require.Equal(t, returnCode, 0)
	})
}
