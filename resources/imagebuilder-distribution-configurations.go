package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ImageBuilderDistributionConfiguration struct {
	svc *imagebuilder.Imagebuilder
	arn string
}

func init() {
	register("ImageBuilderDistributionConfiguration", ListImageBuilderDistributionConfigurations)
}

func ListImageBuilderDistributionConfigurations(sess *session.Session) ([]Resource, error) {
	svc := imagebuilder.New(sess)
	params := &imagebuilder.ListDistributionConfigurationsInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListDistributionConfigurations(params)
		if err != nil {
			return nil, err
		}

		for _, out := range resp.DistributionConfigurationSummaryList {
			resources = append(resources, &ImageBuilderDistributionConfiguration{
				svc: svc,
				arn: *out.Arn,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params = &imagebuilder.ListDistributionConfigurationsInput{
			NextToken: resp.NextToken,
		}
	}

	return resources, nil
}

func (e *ImageBuilderDistributionConfiguration) Remove() error {
	_, err := e.svc.DeleteDistributionConfiguration(&imagebuilder.DeleteDistributionConfigurationInput{
		DistributionConfigurationArn: &e.arn,
	})
	return err
}

func (e *ImageBuilderDistributionConfiguration) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("arn", e.arn)
	return properties
}

func (e *ImageBuilderDistributionConfiguration) String() string {
	return e.arn
}
