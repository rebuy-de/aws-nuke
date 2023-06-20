package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CloudFrontFunction struct {
	svc   *cloudfront.CloudFront
	name  *string
	stage *string
}

func init() {
	register("CloudFrontFunction", ListCloudFrontFunctions)
}

func ListCloudFrontFunctions(sess *session.Session) ([]Resource, error) {
	svc := cloudfront.New(sess)
	resources := []Resource{}

	for {
		resp, err := svc.ListFunctions(nil)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.FunctionList.Items {
			resources = append(resources, &CloudFrontFunction{
				svc:   svc,
				name:  item.Name,
				stage: item.FunctionMetadata.Stage,
			})
		}
		return resources, nil
	}
}

func (f *CloudFrontFunction) Remove() error {
	resp, err := f.svc.GetFunction(&cloudfront.GetFunctionInput{
		Name:  f.name,
		Stage: f.stage,
	})
	if err != nil {
		return err
	}

	_, err = f.svc.DeleteFunction(&cloudfront.DeleteFunctionInput{
		Name:    f.name,
		IfMatch: resp.ETag,
	})

	return err
}

func (f *CloudFrontFunction) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("name", f.name)
	properties.Set("stage", f.stage)
	return properties
}
