package cmd

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type NukeConfig struct {
	AccountBlacklist []string `yaml:"account-blacklist"`
	Region           string   `yaml:"region"`
	Accounts         struct {
		Filters map[string][]string `yaml:"filters"`
	}
}

func LoadConfig(path string) (*NukeConfig, error) {
	var err error

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := new(NukeConfig)
	err = yaml.Unmarshal(raw, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *NukeConfig) HasBlacklist() bool {
	return c.AccountBlacklist != nil && len(c.AccountBlacklist) >= 0
}

func (c *NukeConfig) InBlacklist(searchID string) bool {
	for _, blacklistID := range c.AccountBlacklist {
		if blacklistID == searchID {
			return true
		}
	}

	return false
}
