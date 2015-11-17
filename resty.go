package main

import "github.com/sparkymat/resty/args"

func main() {
	command, commandArgs := args.Parse()

	dispatchCommand(command, commandArgs)
}
