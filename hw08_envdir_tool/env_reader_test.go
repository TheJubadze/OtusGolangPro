package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	dirPath := "./testdata/env/"

	env, err := ReadDir(dirPath)

	require.Nil(t, err)
	require.NotNil(t, env)
	require.Equal(t, EnvValue{Value: "bar", NeedRemove: false}, env["BAR"])
	require.Equal(t, EnvValue{Value: "", NeedRemove: false}, env["EMPTY"])
	require.Equal(t, EnvValue{Value: "   foo\nwith new line", NeedRemove: false}, env["FOO"])
	require.Equal(t, EnvValue{Value: "\"hello\"", NeedRemove: false}, env["HELLO"])
	require.Equal(t, EnvValue{Value: "", NeedRemove: true}, env["UNSET"])
}
