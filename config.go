package main

type NukeConfig struct {
	AccountBlacklist []string `yaml:"account-blacklist"`
	Region           string   `yaml:"region"`
	Accounts         struct {
		Filters map[string][]string `yaml:"filters"`
	}
}
