package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
)

type ConfigServiceConformancePack struct {
	svc                 *configservice.ConfigService
	conformancePackName *string
}

func init() {
	register("ConfigServiceConformancePack", ListConfigServiceConformancePacks)
}

func ListConfigServiceConformancePacks(sess *session.Session) ([]Resource, error) {
	svc := configservice.New(sess)
	resources := []Resource{}

	params := &configservice.DescribeConformancePacksInput{}

	for {
		output, err := svc.DescribeConformancePacks(params)
		if err != nil {
			return nil, err
		}

		for _, conformancePack := range output.ConformancePackDetails {
			resources = append(resources, &ConfigServiceConformancePack{
				svc:                 svc,
				conformancePackName: conformancePack.ConformancePackName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (c *ConfigServiceConformancePack) Remove() error {
	_, err := c.svc.DeleteConformancePack(&configservice.DeleteConformancePackInput{
		ConformancePackName: c.conformancePackName,
	})

	return err
}

func (c *ConfigServiceConformancePack) String() string {
	return *c.conformancePackName
}
