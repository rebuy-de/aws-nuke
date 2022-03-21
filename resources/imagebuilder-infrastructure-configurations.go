package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ImageBuilderInfrastructureConfiguration struct {
	svc *imagebuilder.Imagebuilder
	arn string
}

func init() {
	register("ImageBuilderInfrastructureConfiguration", ListImageBuilderInfrastructureConfigurations)
}

func ListImageBuilderInfrastructureConfigurations(sess *session.Session) ([]Resource, error) {
	svc := imagebuilder.New(sess)
	params := &imagebuilder.ListInfrastructureConfigurationsInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListInfrastructureConfigurations(params)
		if err != nil {
			return nil, err
		}

		for _, out := range resp.InfrastructureConfigurationSummaryList {
			resources = append(resources, &ImageBuilderInfrastructureConfiguration{
				svc: svc,
				arn: *out.Arn,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params = &imagebuilder.ListInfrastructureConfigurationsInput{
			NextToken: resp.NextToken,
		}
	}

	return resources, nil
}

func (e *ImageBuilderInfrastructureConfiguration) Remove() error {
	_, err := e.svc.DeleteInfrastructureConfiguration(&imagebuilder.DeleteInfrastructureConfigurationInput{
		InfrastructureConfigurationArn: &e.arn,
	})
	return err
}

func (e *ImageBuilderInfrastructureConfiguration) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("arn", e.arn)
	return properties
}

func (e *ImageBuilderInfrastructureConfiguration) String() string {
	return e.arn
}
