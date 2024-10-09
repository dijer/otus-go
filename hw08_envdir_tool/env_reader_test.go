package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("test testdata folder", func(t *testing.T) {
		envs, err := ReadDir("./testdata/env")

		require.Nil(t, err)
		require.Equal(t, Environment{
			"BAR":   EnvValue{"bar", false},
			"EMPTY": EnvValue{"", true},
			"FOO":   EnvValue{"   foo\nwith new line", false},
			"HELLO": EnvValue{"\"hello\"", false},
			"UNSET": EnvValue{"", true},
		}, envs)
	})

	t.Run("should skip folder", func(t *testing.T) {
		envs, err := ReadDir("./testdata")

		require.Nil(t, err)
		require.Equal(t, envs, Environment{
			"echo.sh": EnvValue{"#!/usr/bin/env bash", false},
		})
	})
}
