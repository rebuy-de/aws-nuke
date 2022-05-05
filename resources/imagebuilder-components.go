package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ImageBuilderComponent struct {
	svc *imagebuilder.Imagebuilder
	arn string
}

func init() {
	register("ImageBuilderComponent", ListImageBuilderComponents)
}

func ListImageBuilderComponents(sess *session.Session) ([]Resource, error) {
	svc := imagebuilder.New(sess)
	params := &imagebuilder.ListComponentsInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListComponents(params)
		if err != nil {
			return nil, err
		}

		for _, out := range resp.ComponentVersionList {
			resources, err = ListImageBuilderComponentVersions(svc, out.Arn, resources)
			if err != nil {
				return nil, err
			}
		}

		if resp.NextToken == nil {
			break
		}

		params = &imagebuilder.ListComponentsInput{
			NextToken: resp.NextToken,
		}
	}

	return resources, nil
}

func ListImageBuilderComponentVersions(svc *imagebuilder.Imagebuilder, componentVersionArn *string, resources []Resource) ([]Resource, error) {
	params := &imagebuilder.ListComponentBuildVersionsInput{
		ComponentVersionArn: componentVersionArn,
	}

	for {
		resp, err := svc.ListComponentBuildVersions(params)
		if err != nil {
			return nil, err
		}

		for _, out := range resp.ComponentSummaryList {
			resources = append(resources, &ImageBuilderComponent{
				svc: svc,
				arn: *out.Arn,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params = &imagebuilder.ListComponentBuildVersionsInput{
			ComponentVersionArn: componentVersionArn,
			NextToken:           resp.NextToken,
		}
	}
	return resources, nil
}

func (e *ImageBuilderComponent) Remove() error {
	_, err := e.svc.DeleteComponent(&imagebuilder.DeleteComponentInput{
		ComponentBuildVersionArn: &e.arn,
	})
	return err
}

func (e *ImageBuilderComponent) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("arn", e.arn)
	return properties
}

func (e *ImageBuilderComponent) String() string {
	return e.arn
}
