package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type GlueSecurityConfiguration struct {
	svc  *glue.Glue
	name *string
}

func init() {
	register("GlueSecurityConfiguration", ListGlueSecurityConfigurations)
}

func ListGlueSecurityConfigurations(sess *session.Session) ([]Resource, error) {
	svc := glue.New(sess)
	resources := []Resource{}

	params := &glue.GetSecurityConfigurationsInput{
		MaxResults: aws.Int64(25),
	}

	for {
		output, err := svc.GetSecurityConfigurations(params)
		if err != nil {
			return nil, err
		}

		for _, securityConfiguration := range output.SecurityConfigurations {
			resources = append(resources, &GlueSecurityConfiguration{
				svc:  svc,
				name: securityConfiguration.Name,
			})
		}

		// Check if there are more security configurations to fetch
		if output.NextToken == nil || *output.NextToken == "" {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueSecurityConfiguration) Remove() error {
	_, err := f.svc.DeleteSecurityConfiguration(&glue.DeleteSecurityConfigurationInput{
		Name: f.name,
	})

	return err
}

func (f *GlueSecurityConfiguration) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", f.name)

	return properties
}

func (f *GlueSecurityConfiguration) String() string {
	return *f.name
}
