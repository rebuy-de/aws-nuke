package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesisvideo"
)

type KinesisSignalingChannels struct {
	svc        *kinesisvideo.KinesisVideo
	ChannelARN *string
}

func init() {
	register("KinesisSignalingChannels", ListKinesisSignalingChannels)
}

func ListKinesisSignalingChannels(sess *session.Session) ([]Resource, error) {
	svc := kinesisvideo.New(sess)
	resources := []Resource{}

	params := &kinesisvideo.ListSignalingChannelsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListSignalingChannels(params)
		if err != nil {
			return nil, err
		}

		for _, streamInfo := range output.ChannelInfoList {
			resources = append(resources, &KinesisSignalingChannels{
				svc:        svc,
				ChannelARN: streamInfo.ChannelARN,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *KinesisSignalingChannels) Remove() error {

	_, err := f.svc.DeleteSignalingChannel(&kinesisvideo.DeleteSignalingChannelInput{
		ChannelARN: f.ChannelARN,
	})

	return err
}

func (f *KinesisSignalingChannels) String() string {
	return *f.ChannelARN
}
