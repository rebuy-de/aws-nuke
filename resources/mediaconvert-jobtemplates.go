package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediaconvert"
)

type MediaConvertJobTemplate struct {
	svc  *mediaconvert.MediaConvert
	name *string
}

func init() {
	register("MediaConvertJobTemplate", ListMediaConvertJobTemplates)
}

func ListMediaConvertJobTemplates(sess *session.Session) ([]Resource, error) {
	svc := mediaconvert.New(sess)
	resources := []Resource{}
	var mediaEndpoint *string

	output, err := svc.DescribeEndpoints(&mediaconvert.DescribeEndpointsInput{})
	if err != nil {
		return nil, err
	}

	for _, mediaconvert := range output.Endpoints {
		mediaEndpoint = mediaconvert.Url
	}

	// Update svc to use custom media endpoint
	svc.Endpoint = *mediaEndpoint

	params := &mediaconvert.ListJobTemplatesInput{
		MaxResults: aws.Int64(20),
	}

	for {
		output, err := svc.ListJobTemplates(params)
		if err != nil {
			return nil, err
		}

		for _, jobTemplate := range output.JobTemplates {
			resources = append(resources, &MediaConvertJobTemplate{
				svc:  svc,
				name: jobTemplate.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *MediaConvertJobTemplate) Remove() error {

	_, err := f.svc.DeleteJobTemplate(&mediaconvert.DeleteJobTemplateInput{
		Name: f.name,
	})

	return err
}

func (f *MediaConvertJobTemplate) String() string {
	return *f.name
}
