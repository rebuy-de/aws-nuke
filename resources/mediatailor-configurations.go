package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediatailor"
)

type MediaTailorConfiguration struct {
	svc  *mediatailor.MediaTailor
	name *string
}

func init() {
	register("MediaTailorConfiguration", ListMediaTailorConfigurations)
}

func ListMediaTailorConfigurations(sess *session.Session) ([]Resource, error) {
	svc := mediatailor.New(sess)
	resources := []Resource{}

	params := &mediatailor.ListPlaybackConfigurationsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		resp, err := svc.ListPlaybackConfigurations(params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {
			resources = append(resources, &MediaTailorConfiguration{
				svc:  svc,
				name: item.Name,
			})
		}
		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}
	return resources, nil
}

func (f *MediaTailorConfiguration) Remove() error {

	_, err := f.svc.DeletePlaybackConfiguration(&mediatailor.DeletePlaybackConfigurationInput{
		Name: f.name,
	})

	return err
}

func (f *MediaTailorConfiguration) String() string {
	return *f.name
}
