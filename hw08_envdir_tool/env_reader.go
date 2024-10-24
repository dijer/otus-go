package main

import (
	"bufio"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
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
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envs := make(Environment)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			return nil, err
		}

		if strings.Contains(fileInfo.Name(), "=") {
			continue
		}

		if strings.Contains(fileInfo.Name(), " ") {
			continue
		}

		if fileInfo.Size() == 0 {
			envs[fileInfo.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		envValue, err := getFileEnv(dir, file)
		if err != nil {
			return nil, err
		}

		envs[fileInfo.Name()] = *envValue
	}

	return envs, nil
}

func getFileEnv(dir string, entry fs.DirEntry) (*EnvValue, error) {
	fileInfo, err := entry.Info()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filepath.Join(dir, fileInfo.Name()))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	firstLine, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	envValue := strings.TrimRight(firstLine, " \t\n")
	envValue = strings.ReplaceAll(envValue, "\x00", "\n")

	if envValue == "" {
		return &EnvValue{NeedRemove: true}, nil
	}

	return &EnvValue{Value: envValue}, nil
}
