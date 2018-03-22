package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediaconvert"
)

type MediaConvertPreset struct {
	svc  *mediaconvert.MediaConvert
	name *string
}

func init() {
	register("MediaConvertPreset", ListMediaConvertPresets)
}

func ListMediaConvertPresets(sess *session.Session) ([]Resource, error) {
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

	params := &mediaconvert.ListPresetsInput{
		MaxResults: aws.Int64(20),
	}

	for {
		output, err := svc.ListPresets(params)
		if err != nil {
			return nil, err
		}

		for _, preset := range output.Presets {
			resources = append(resources, &MediaConvertPreset{
				svc:  svc,
				name: preset.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *MediaConvertPreset) Remove() error {

	_, err := f.svc.DeletePreset(&mediaconvert.DeletePresetInput{
		Name: f.name,
	})

	return err
}

func (f *MediaConvertPreset) String() string {
	return *f.name
}
