package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type SSMParameter struct {
	svc  *ssm.SSM
	name *string
	tags []*ssm.Tag
}

func init() {
	register("SSMParameter", ListSSMParameters)
}

func ListSSMParameters(sess *session.Session) ([]Resource, error) {
	svc := ssm.New(sess)
	resources := []Resource{}

	params := &ssm.DescribeParametersInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.DescribeParameters(params)
		if err != nil {
			return nil, err
		}

		for _, parameter := range output.Parameters {
			tagParams := &ssm.ListTagsForResourceInput{
				ResourceId:   parameter.Name,
				ResourceType: aws.String(ssm.ResourceTypeForTaggingParameter),
			}

			tagResp, tagErr := svc.ListTagsForResource(tagParams)
			if tagErr != nil {
				return nil, tagErr
			}

			resources = append(resources, &SSMParameter{
				svc:  svc,
				name: parameter.Name,
				tags: tagResp.TagList,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *SSMParameter) Remove() error {

	_, err := f.svc.DeleteParameter(&ssm.DeleteParameterInput{
		Name: f.name,
	})

	return err
}

func (f *SSMParameter) String() string {
	return *f.name
}

func (f *SSMParameter) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tag := range f.tags {
		properties.SetTag(tag.Key, tag.Value)
	}
	properties.
		Set("Name", f.name)
	return properties
}
