package cmd

import (
	"fmt"
	"strings"
)

type NukeParameters struct {
	ConfigPath string

	Profile         string
	AccessKeyID     string
	SecretAccessKey string

	NoDryRun bool
	Force    bool
}

func (p *NukeParameters) Validate() error {
	if strings.TrimSpace(p.ConfigPath) == "" {
		return fmt.Errorf("You have to specify the --config flag.\n")
	}

	if p.hasProfile() == p.hasKeys() {
		return fmt.Errorf("You have to specify the --profile flag OR " +
			"--access-key-id and --secret-access-key.\n")
	}

	return nil
}

func (p *NukeParameters) hasProfile() bool {
	return strings.TrimSpace(p.Profile) != ""
}

func (p *NukeParameters) hasKeys() bool {
	return strings.TrimSpace(p.AccessKeyID) != "" &&
		strings.TrimSpace(p.SecretAccessKey) != ""
}
