package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Environment map[string]string

func processBytesValue(value []byte) []byte {
	value = bytes.TrimRight(value, " \t\n")
	value = bytes.ReplaceAll(value, []byte("\x00"), []byte("\n"))
	return value
}

func readFirstLineOfFile(file *os.File) ([]byte, error) {
	reader := bufio.NewReader(file)
	value, err := reader.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("error while reading first line of %s: %w", file.Name(), err)
	}
	return value, nil
}

func getValueFromFile(dir, fileName string) (string, error) {
	if strings.Contains(fileName, "=") {
		return "", fmt.Errorf("failed to process %s: it's name contains '='", fileName)
	}
	filePath := path.Join(dir, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file in dir: %w", err)
	}
	defer file.Close()
	value, err := readFirstLineOfFile(file)
	if err != nil {
		return "", err
	}
	value = processBytesValue(value)

	return string(value), nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	fileNames, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory data: %w", err)
	}

	env := make(Environment, len(fileNames))
	for _, fileName := range fileNames {
		value, err := getValueFromFile(dir, fileName.Name())
		if err != nil {
			return nil, err
		}
		env[fileName.Name()] = value
	}

	return env, nil
}
