package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		log.Print("No command provided")
		return 1
	}

	if !isSafeCommand(cmd) {
		log.Print("Unsafe command detected")
		return 1
	}

	envSet := mapEnvironmentToStringSlice(env)
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = append(os.Environ(), envSet...)

	err := command.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		log.Print("Command execution failed: ", err)
		return 1
	}

	return 0
}

func isSafeCommand(cmd []string) bool {
	for _, arg := range cmd {
		if strings.ContainsAny(arg, "&;|") {
			return false
		}
	}
	return true
}

func mapEnvironmentToStringSlice(env Environment) []string {
	result := make([]string, 0)

	for key, value := range env {
		if value.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				log.Print(err)
			}
			continue
		}
		envString := fmt.Sprintf("%s=%s", key, value.Value)
		result = append(result, envString)
	}

	return result
}
