package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/medialive"
)

type MediaLiveInput struct {
	svc *medialive.MediaLive
	ID  *string
}

func init() {
	register("MediaLiveInput", ListMediaLiveInputs)
}

func ListMediaLiveInputs(sess *session.Session) ([]Resource, error) {
	svc := medialive.New(sess)
	resources := []Resource{}

	params := &medialive.ListInputsInput{
		MaxResults: aws.Int64(20),
	}

	for {
		output, err := svc.ListInputs(params)
		if err != nil {
			return nil, err
		}

		for _, input := range output.Inputs {
			resources = append(resources, &MediaLiveInput{
				svc: svc,
				ID:  input.Id,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *MediaLiveInput) Remove() error {

	_, err := f.svc.DeleteInput(&medialive.DeleteInputInput{
		InputId: f.ID,
	})

	return err
}

func (f *MediaLiveInput) String() string {
	return *f.ID
}
