package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/medialive"
)

type MediaLiveChannel struct {
	svc *medialive.MediaLive
	ID  *string
}

func init() {
	register("MediaLiveChannel", ListMediaLiveChannels)
}

func ListMediaLiveChannels(sess *session.Session) ([]Resource, error) {
	svc := medialive.New(sess)
	resources := []Resource{}

	params := &medialive.ListChannelsInput{
		MaxResults: aws.Int64(20),
	}

	for {
		output, err := svc.ListChannels(params)
		if err != nil {
			return nil, err
		}

		for _, channel := range output.Channels {
			resources = append(resources, &MediaLiveChannel{
				svc: svc,
				ID:  channel.Id,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *MediaLiveChannel) Remove() error {

	_, err := f.svc.DeleteChannel(&medialive.DeleteChannelInput{
		ChannelId: f.ID,
	})

	return err
}

func (f *MediaLiveChannel) String() string {
	return *f.ID
}
