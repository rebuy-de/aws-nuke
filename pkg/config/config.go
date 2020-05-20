package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/rebuy-de/aws-nuke/pkg/types"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ResourceTypes struct {
	Targets  types.Collection `yaml:"targets"`
	Excludes types.Collection `yaml:"excludes"`
}

type Account struct {
	Filters       Filters       `yaml:"filters"`
	ResourceTypes ResourceTypes `yaml:"resource-types"`
	Presets       []string      `yaml:"presets"`
}

type Nuke struct {
	AccountBlacklist []string                     `yaml:"account-blacklist"`
	Regions          []string                     `yaml:"regions"`
	Accounts         map[string]Account           `yaml:"accounts"`
	ResourceTypes    ResourceTypes                `yaml:"resource-types"`
	Presets          map[string]PresetDefinitions `yaml:"presets"`
	FeatureFlags     FeatureFlags                 `yaml:"feature-flags"`
	CustomEndpoints  CustomEndpoints              `yaml:"endpoints"`
}

type FeatureFlags struct {
	DisableDeletionProtection struct {
		RDSInstance         bool `yaml:"RDSInstance"`
		EC2Instance         bool `yaml:"EC2Instance"`
		CloudformationStack bool `yaml:"CloudformationStack"`
	} `yaml:"disable-deletion-protection"`
}

type PresetDefinitions struct {
	Filters Filters `yaml:"filters"`
}

type CustomService struct {
	Service               string `yaml:"service"`
	URL                   string `yaml:"url"`
	TLSInsecureSkipVerify bool   `yaml:"tls_insecure_skip_verify"`
}

type CustomServices []*CustomService

type CustomRegion struct {
	Region                string         `yaml:"region"`
	Services              CustomServices `yaml:"services"`
	TLSInsecureSkipVerify bool           `yaml:"tls_insecure_skip_verify"`
}

type CustomEndpoints []*CustomRegion

func Load(path string) (*Nuke, error) {
	var err error

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := new(Nuke)
	err = yaml.UnmarshalStrict(raw, config)
	if err != nil {
		return nil, err
	}

	if err := config.resolveDeprecations(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Nuke) HasBlacklist() bool {
	return c.AccountBlacklist != nil && len(c.AccountBlacklist) > 0
}

func (c *Nuke) InBlacklist(searchID string) bool {
	for _, blacklistID := range c.AccountBlacklist {
		if blacklistID == searchID {
			return true
		}
	}

	return false
}

func (c *Nuke) ValidateAccount(accountID string, aliases []string) error {
	if !c.HasBlacklist() {
		return fmt.Errorf("The config file contains an empty blacklist. " +
			"For safety reasons you need to specify at least one account ID. " +
			"This should be your production account.")
	}

	if c.InBlacklist(accountID) {
		return fmt.Errorf("You are trying to nuke the account with the ID %s, "+
			"but it is blacklisted. Aborting.", accountID)
	}

	if len(aliases) == 0 {
		return fmt.Errorf("The specified account doesn't have an alias. " +
			"For safety reasons you need to specify an account alias. " +
			"Your production account should contain the term 'prod'.")
	}

	for _, alias := range aliases {
		if strings.Contains(strings.ToLower(alias), "prod") {
			return fmt.Errorf("You are trying to nuke an account with the alias '%s', "+
				"but it has the substring 'prod' in it. Aborting.", alias)
		}
	}

	if _, ok := c.Accounts[accountID]; !ok {
		return fmt.Errorf("Your account ID '%s' isn't listed in the config. "+
			"Aborting.", accountID)
	}

	return nil
}

func (c *Nuke) Filters(accountID string) (Filters, error) {
	account := c.Accounts[accountID]
	filters := account.Filters

	if filters == nil {
		filters = Filters{}
	}

	if account.Presets == nil {
		return filters, nil
	}

	for _, presetName := range account.Presets {
		notFound := fmt.Errorf("Could not find filter preset '%s'", presetName)
		if c.Presets == nil {
			return nil, notFound
		}

		preset, ok := c.Presets[presetName]
		if !ok {
			return nil, notFound
		}

		filters.Merge(preset.Filters)
	}

	return filters, nil
}

func (c *Nuke) resolveDeprecations() error {
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
			log.Warnf("deprecated resource type '%s' - converting to '%s'\n", resourceType, replacement)

			if _, ok := a.Filters[replacement]; ok {
				return fmt.Errorf("using deprecated resource type and replacement: '%s','%s'", resourceType, replacement)
			}

			a.Filters[replacement] = resources
			delete(a.Filters, resourceType)
		}
	}
	return nil
}

// GetRegion returns the custom region or nil when no such custom endpoints are defined for this region
func (endpoints CustomEndpoints) GetRegion(region string) *CustomRegion {
	for _, r := range endpoints {
		if r.Region == region {
			if r.TLSInsecureSkipVerify {
				for _, s := range r.Services {
					s.TLSInsecureSkipVerify = r.TLSInsecureSkipVerify
				}
			}
			return r
		}
	}
	return nil
}

// GetService returns the custom region or nil when no such custom endpoints are defined for this region
func (services CustomServices) GetService(serviceType string) *CustomService {
	for _, s := range services {
		if serviceType == s.Service {
			return s
		}
	}
	return nil
}

func (endpoints CustomEndpoints) GetURL(region, serviceType string) string {
	r := endpoints.GetRegion(region)
	if r == nil {
		return ""
	}
	s := r.Services.GetService(serviceType)
	if s == nil {
		return ""
	}
	return s.URL
}
