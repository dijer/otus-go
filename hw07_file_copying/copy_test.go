package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrFromPathRequired(t *testing.T) {
	err := Copy("", "./testdata/output.txt", 0, 0)
	require.ErrorIs(t, err, ErrFromPathRequired)
}

func TestErrToPathRequired(t *testing.T) {
	err := Copy("./testdata/input.txt", "", 0, 0)
	require.ErrorIs(t, err, ErrToPathRequired)
}

func TestErrSameFile(t *testing.T) {
	err := Copy("./testdata/input.txt", "./testdata/input.txt", 0, 0)
	require.ErrorIs(t, err, ErrSameFile)
}

func TestErrLimitLessZero(t *testing.T) {
	err := Copy("./testdata/input.txt", "./testdata/output.txt", 0, -10)
	require.ErrorIs(t, err, ErrLimitLessZero)
}

func TestErrOffsetLessZero(t *testing.T) {
	err := Copy("./testdata/input.txt", "./testdata/output.txt", -10, 0)
	require.ErrorIs(t, err, ErrOffsetLessZero)
}

func TestErrFileNotFound(t *testing.T) {
	err := Copy("./testdata/somefile.txt", "./testdata/output.txt", 0, 0)
	require.ErrorIs(t, err, ErrFileNotFound)
}
