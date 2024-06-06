package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	pflag.Parse()
	args := pflag.Args()

	if len(args) < 2 {
		_, err := fmt.Fprintln(os.Stderr, "usage: go-envdir dir child")
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}

	env, err := ReadDir(args[0])
	if err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}

	os.Exit(RunCmd(args[1:], env))
}
