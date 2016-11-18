package main

import (
	"fmt"
	"os"

	flag "github.com/ogier/pflag"
)

var (
	// will be overwritten on build
	version = "unknown"
)

type NukeParameters struct {
	ConfigPath string

	Profile         string
	AccessKeyID     string
	SecretAccessKey string

	NoDryRun bool
	Force    bool
}

func main() {
	var err error

	fmt.Printf("Running aws-nuke version %s.\n", version)

	params := NukeParameters{}

	flag.StringVar(&params.ConfigPath, "config", "", "path to config")
	flag.StringVar(&params.Profile, "aws-profile", "", "profile name to nuke")
	flag.StringVar(&params.AccessKeyID, "aws-access-key-id", "", "AWS access-key-id")
	flag.StringVar(&params.SecretAccessKey, "secret-access-key", "", "AWS secret-access-key")
	flag.BoolVar(&params.NoDryRun, "no-dry-run", false, "Actualy delete found resources.")
	flag.BoolVar(&params.Force, "force", false, "Don't ask for confirmation.")

	flag.Parse()

	if !params.NoDryRun {
		fmt.Printf("Dry run: do real delete with '--no-dry-run'.\n")
	}

	fmt.Println()

	n := NewNuke(params)

	err = n.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = n.StartSession()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	n.Run()
}
