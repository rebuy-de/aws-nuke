package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ImageBuilderPipeline struct {
	svc *imagebuilder.Imagebuilder
	arn string
}

func init() {
	register("ImageBuilderPipeline", ListImageBuilderPipelines)
}

func ListImageBuilderPipelines(sess *session.Session) ([]Resource, error) {
	svc := imagebuilder.New(sess)
	params := &imagebuilder.ListImagePipelinesInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListImagePipelines(params)

		if err != nil {
			return nil, err
		}

		for _, out := range resp.ImagePipelineList {
			resources = append(resources, &ImageBuilderPipeline{
				svc: svc,
				arn: *out.Arn,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params = &imagebuilder.ListImagePipelinesInput{
			NextToken: resp.NextToken,
		}
	}
	return resources, nil
}

func (e *ImageBuilderPipeline) Remove() error {
	_, err := e.svc.DeleteImagePipeline(&imagebuilder.DeleteImagePipelineInput{
		ImagePipelineArn: &e.arn,
	})
	return err
}

func (e *ImageBuilderPipeline) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("arn", e.arn)
	return properties
}

func (e *ImageBuilderPipeline) String() string {
	return e.arn
}
