package main

import (
	"os"

	"github.com/rebuy-de/aws-nuke/cmd"
)

// NukeParameters defines the command-line parameters for aws-nuke.
type NukeParameters struct {
	ConfigPath string

	Profile         string
	AccessKeyID     string
	SecretAccessKey string

	NoDryRun   bool
	Force      bool
	ForceSleep int
	Quiet      bool

	MaxWaitRetries int
}

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		os.Exit(-1)
	}
}
