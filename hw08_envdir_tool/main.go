package main

import (
	"log"
	"os"
)

func main() {
	programName := os.Args[0]

	if len(os.Args) < 3 {
		log.Fatalf(
			"Program shoud be called as: %s <env files directory> <command> <optional arguments>",
			programName)
	}

	envDir := os.Args[1]
	commandAndArguments := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		log.Fatalf("Unable to get environment variables info: %v", err)
	}
	code := RunCmd(commandAndArguments, env)
	os.Exit(code)
}
