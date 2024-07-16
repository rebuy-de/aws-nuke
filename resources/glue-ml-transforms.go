package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type GlueMLTransform struct {
	svc *glue.Glue
	id  *string
}

func init() {
	register("GlueMLTransform", ListGlueMLTransforms)
}

func ListGlueMLTransforms(sess *session.Session) ([]Resource, error) {
	svc := glue.New(sess)
	resources := []Resource{}

	params := &glue.ListMLTransformsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListMLTransforms(params)
		if err != nil {
			return nil, err
		}

		for _, transformId := range output.TransformIds {
			resources = append(resources, &GlueMLTransform{
				svc: svc,
				id:  transformId,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueMLTransform) Remove() error {
	_, err := f.svc.DeleteMLTransform(&glue.DeleteMLTransformInput{
		TransformId: f.id,
	})

	return err
}

func (f *GlueMLTransform) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Id", f.id)

	return properties
}

func (f *GlueMLTransform) String() string {
	return *f.id
}
