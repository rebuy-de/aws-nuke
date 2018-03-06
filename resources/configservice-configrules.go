package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
)

type ConfigServiceConfigRule struct {
	svc            *configservice.ConfigService
	configRuleName *string
}

func init() {
	register("ConfigServiceConfigRule", ListConfigServiceConfigRules)
}

func ListConfigServiceConfigRules(sess *session.Session) ([]Resource, error) {
	svc := configservice.New(sess)
	resources := []Resource{}

	params := &configservice.DescribeConfigRulesInput{}

	for {
		output, err := svc.DescribeConfigRules(params)
		if err != nil {
			return nil, err
		}

		for _, configRule := range output.ConfigRules {
			resources = append(resources, &ConfigServiceConfigRule{
				svc:            svc,
				configRuleName: configRule.ConfigRuleName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *ConfigServiceConfigRule) Remove() error {

	_, err := f.svc.DeleteConfigRule(&configservice.DeleteConfigRuleInput{
		ConfigRuleName: f.configRuleName,
	})

	return err
}

func (f *ConfigServiceConfigRule) String() string {
	return *f.configRuleName
}
