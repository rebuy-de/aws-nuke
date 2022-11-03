package main

import (
	"os"

	"github.com/hunterkepley/aws-nuke/v2/cmd"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		os.Exit(-1)
	}
}
