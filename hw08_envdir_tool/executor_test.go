package main

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmdPositive(t *testing.T) {
	var cmd []string
	switch runtime.GOOS {
	case "windows":
		cmd = []string{"cmd", "-c", "echo \"hello world\""}
	default:
		cmd = []string{"sh", "-c", "echo \"hello world\""}
	}
	env := make(Environment)
	env["BAR"] = EnvValue{Value: "bar"}

	errorCode := RunCmd(cmd, env)

	require.Equal(t, 0, errorCode)
}

func TestRunCmdNegative(t *testing.T) {
	cmd := []string{"unknown_command", "-c", "echo \"hello world\""}
	env := make(Environment)
	env["BAR"] = EnvValue{Value: "bar"}

	errorCode := RunCmd(cmd, env)

	require.Equal(t, 1, errorCode)
}
