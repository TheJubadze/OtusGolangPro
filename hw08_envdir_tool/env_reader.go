package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
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
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment)
	for _, entry := range dirEntries {
		if !entry.IsDir() && !strings.Contains(entry.Name(), "=") {
			env[entry.Name()], err = readEntry(filepath.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}
		}
	}

	return env, nil
}

func readEntry(filePath string) (EnvValue, error) {
	envFile, err := os.Open(filePath)
	if err != nil {
		return EnvValue{}, err
	}
	defer closeFile(envFile)

	fi, err := envFile.Stat()
	if err != nil {
		return EnvValue{}, err
	}

	if fi.Size() == 0 {
		return EnvValue{"", true}, nil
	}

	reader := bufio.NewReader(envFile)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return EnvValue{}, err
	}

	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")
	line = strings.TrimRight(line, " \t")
	processedLine := bytes.ReplaceAll([]byte(line), []byte{0x00}, []byte{'\n'})

	return EnvValue{string(processedLine), false}, nil
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Println(err)
	}
}
