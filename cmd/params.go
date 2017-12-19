package cmd

import (
	"fmt"
	"strings"
)

type NukeParameters struct {
	ConfigPath string

	Targets []string

	NoDryRun bool
	Force    bool
}

func (p *NukeParameters) Validate() error {
	if strings.TrimSpace(p.ConfigPath) == "" {
		return fmt.Errorf("You have to specify the --config flag.\n")
	}

	return nil
}

func (p *NukeParameters) WantsTarget(name string) bool {
	if p.Targets == nil || len(p.Targets) < 1 {
		return true
	}

	for _, wants := range p.Targets {
		if wants == name {
			return true
		}
	}

	return false
}
