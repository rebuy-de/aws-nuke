package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediapackage"
)

type MediaPackageChannel struct {
	svc *mediapackage.MediaPackage
	ID  *string
}

func init() {
	register("MediaPackageChannel", ListMediaPackageChannels)
}

func ListMediaPackageChannels(sess *session.Session) ([]Resource, error) {
	svc := mediapackage.New(sess)
	resources := []Resource{}

	params := &mediapackage.ListChannelsInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.ListChannels(params)
		if err != nil {
			return nil, err
		}

		for _, channel := range output.Channels {
			resources = append(resources, &MediaPackageChannel{
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

func (f *MediaPackageChannel) Remove() error {

	_, err := f.svc.DeleteChannel(&mediapackage.DeleteChannelInput{
		Id: f.ID,
	})

	return err
}

func (f *MediaPackageChannel) String() string {
	return *f.ID
}
