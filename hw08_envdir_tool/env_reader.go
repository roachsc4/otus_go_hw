package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Environment map[string]string

func processValue(value string) string {
	value = strings.TrimRight(value, " \t\n")
	value = strings.ReplaceAll(value, "\x00", "\n")
	return value
}

func readFirstLineOfFile(file *os.File) (string, error) {
	reader := bufio.NewReader(file)
	value, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("error while reading first line of %s: %w", file.Name(), err)
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
	value = processValue(value)

	return value, nil
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
