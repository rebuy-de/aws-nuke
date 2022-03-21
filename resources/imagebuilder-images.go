package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ImageBuilderImage struct {
	svc *imagebuilder.Imagebuilder
	arn string
}

func init() {
	register("ImageBuilderImage", ListImageBuilderImages)
}

func ListImageBuilderImages(sess *session.Session) ([]Resource, error) {
	svc := imagebuilder.New(sess)
	params := &imagebuilder.ListImagesInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListImages(params)
		if err != nil {
			return nil, err
		}

		for _, out := range resp.ImageVersionList {
			resources, err = ImageBuildVersions(svc, out.Arn, resources)
			if err != nil {
				return nil, err
			}
		}

		if resp.NextToken == nil {
			break
		}

		params = &imagebuilder.ListImagesInput{
			NextToken: resp.NextToken,
		}
	}

	return resources, nil
}

func ImageBuildVersions(svc *imagebuilder.Imagebuilder, imageVersionArn *string, resources []Resource) ([]Resource, error) {
	params := &imagebuilder.ListImageBuildVersionsInput{
		ImageVersionArn: imageVersionArn,
	}

	for {
		resp, err := svc.ListImageBuildVersions(params)
		if err != nil {
			return nil, err
		}

		for _, out := range resp.ImageSummaryList {
			resources = append(resources, &ImageBuilderImage{
				svc: svc,
				arn: *out.Arn,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params = &imagebuilder.ListImageBuildVersionsInput{
			ImageVersionArn: imageVersionArn,
			NextToken:       resp.NextToken,
		}
	}
	return resources, nil
}

func (e *ImageBuilderImage) Remove() error {
	_, err := e.svc.DeleteImage(&imagebuilder.DeleteImageInput{
		ImageBuildVersionArn: &e.arn,
	})
	return err
}

func (e *ImageBuilderImage) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("arn", e.arn)
	return properties
}

func (e *ImageBuilderImage) String() string {
	return e.arn
}
