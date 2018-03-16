package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/medialive"
)

type MediaLiveInputSecurityGroup struct {
	svc *medialive.MediaLive
	ID  *string
}

func init() {
	register("MediaLiveInputSecurityGroup", ListMediaLiveInputSecurityGroups)
}

func ListMediaLiveInputSecurityGroups(sess *session.Session) ([]Resource, error) {
	svc := medialive.New(sess)
	resources := []Resource{}

	params := &medialive.ListInputSecurityGroupsInput{
		MaxResults: aws.Int64(20),
	}

	for {
		output, err := svc.ListInputSecurityGroups(params)
		if err != nil {
			return nil, err
		}

		for _, inputSecurityGroup := range output.InputSecurityGroups {
			resources = append(resources, &MediaLiveInputSecurityGroup{
				svc: svc,
				ID:  inputSecurityGroup.Id,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *MediaLiveInputSecurityGroup) Remove() error {

	_, err := f.svc.DeleteInputSecurityGroup(&medialive.DeleteInputSecurityGroupInput{
		InputSecurityGroupId: f.ID,
	})

	return err
}

func (f *MediaLiveInputSecurityGroup) String() string {
	return *f.ID
}
