package cmd

import (
	"fmt"
	"strings"
)

type NukeParameters struct {
	ConfigPath string

	Targets []string
	Include []string
	Exclude []string

	NoDryRun bool
	Force    bool
}

func (p *NukeParameters) Validate() error {
	if strings.TrimSpace(p.ConfigPath) == "" {
		return fmt.Errorf("You have to specify the --config flag.\n")
	}

	if len(p.Targets) > 0 {
		LogWarn("The flag --target is deprecated. Please use --include instead.\n")
	}

	if len(p.Targets) > 0 && len(p.Include) > 0 {
		return fmt.Errorf("The flag --include cannot used together with --target.")
	}

	return nil
}
