package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesisvideo"
)

type KinesisVideoProject struct {
	svc       *kinesisvideo.KinesisVideo
	streamARN *string
}

func init() {
	register("KinesisVideoProject", ListKinesisVideoProjects)
}

func ListKinesisVideoProjects(sess *session.Session) ([]Resource, error) {
	svc := kinesisvideo.New(sess)
	resources := []Resource{}

	params := &kinesisvideo.ListStreamsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListStreams(params)
		if err != nil {
			return nil, err
		}

		for _, streamInfo := range output.StreamInfoList {
			resources = append(resources, &KinesisVideoProject{
				svc:       svc,
				streamARN: streamInfo.StreamARN,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *KinesisVideoProject) Remove() error {

	_, err := f.svc.DeleteStream(&kinesisvideo.DeleteStreamInput{
		StreamARN: f.streamARN,
	})

	return err
}

func (f *KinesisVideoProject) String() string {
	return *f.streamARN
}
