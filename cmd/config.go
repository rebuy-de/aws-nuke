package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type NukeConfig struct {
	AccountBlacklist []string                     `yaml:"account-blacklist"`
	Regions          []string                     `yaml:"regions"`
	Accounts         map[string]NukeConfigAccount `yaml:"accounts"`
}

type NukeConfigAccount struct {
	Filters map[string][]string `yaml:"filters"`
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

	if err := config.resolveDeprecations(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *NukeConfig) HasBlacklist() bool {
	return c.AccountBlacklist != nil && len(c.AccountBlacklist) > 0
}

func (c *NukeConfig) InBlacklist(searchID string) bool {
	for _, blacklistID := range c.AccountBlacklist {
		if blacklistID == searchID {
			return true
		}
	}

	return false
}

func (c *NukeConfig) ValidateAccount(accountID string, aliases []string) (*NukeConfigAccount, error) {
	if !c.HasBlacklist() {
		return nil, fmt.Errorf("The config file contains an empty blacklist. " +
			"For safety reasons you need to specify at least one account ID. " +
			"This should be your production account.")
	}

	if c.InBlacklist(accountID) {
		return nil, fmt.Errorf("You are trying to nuke the account with the ID %s, "+
			"but it is blacklisted. Aborting.", accountID)
	}

	if len(aliases) == 0 {
		return nil, fmt.Errorf("The specified account doesn't have an alias. " +
			"For safety reasons you need to specify an account alias. " +
			"Your production account should contain the term 'prod'.")
	}

	for _, alias := range aliases {
		if strings.Contains(strings.ToLower(alias), "prod") {
			return nil, fmt.Errorf("You are trying to nuke an account with the alias '%s', "+
				"but it has the substring 'prod' in it. Aborting.", alias)
		}
	}

	if _, ok := c.Accounts[accountID]; !ok {
		return nil, fmt.Errorf("Your account ID '%s' isn't listed in the config. "+
			"Aborting.", accountID)
	}

	ac := c.Accounts[accountID]
	return &ac, nil
}

func (c *NukeConfig) resolveDeprecations() error {
	deprecations := map[string]string{
		"EC2DhcpOptions":                "EC2DHCPOptions",
		"EC2InternetGatewayAttachement": "EC2InternetGatewayAttachment",
		"EC2NatGateway":                 "EC2NATGateway",
		"EC2Vpc":                        "EC2VPC",
		"EC2VpcEndpoint":                "EC2VPCEndpoint",
		"EC2VpnConnection":              "EC2VPNConnection",
		"EC2VpnGateway":                 "EC2VPNGateway",
		"EC2VpnGatewayAttachement":      "EC2VPNGatewayAttachment",
		"ECRrepository":                 "ECRRepository",
		"IamGroup":                      "IAMGroup",
		"IamGroupPolicyAttachement":     "IAMGroupPolicyAttachment",
		"IamInstanceProfile":            "IAMInstanceProfile",
		"IamInstanceProfileRole":        "IAMInstanceProfileRole",
		"IamPolicy":                     "IAMPolicy",
		"IamRole":                       "IAMRole",
		"IamRolePolicyAttachement":      "IAMRolePolicyAttachment",
		"IamServerCertificate":          "IAMServerCertificate",
		"IamUser":                       "IAMUser",
		"IamUserAccessKeys":             "IAMUserAccessKey",
		"IamUserGroupAttachement":       "IAMUserGroupAttachment",
		"IamUserPolicyAttachement":      "IAMUserPolicyAttachment",
		"RDSCluster":                    "RDSDBCluster",
	}

	for _, a := range c.Accounts {
		for resourceType, resources := range a.Filters {
			replacement, ok := deprecations[resourceType]
			if !ok {
				continue
			}
			LogWarn("deprecated resource type '%s' - converting to '%s'\n", resourceType, replacement)

			if _, ok := a.Filters[replacement]; ok {
				return fmt.Errorf("using deprecated resource type and replacement: '%s','%s'", resourceType, replacement)
			}

			a.Filters[replacement] = resources
			delete(a.Filters, resourceType)
		}
	}
	return nil
}
