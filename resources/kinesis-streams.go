package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

type KinesisStream struct {
	svc        *kinesis.Kinesis
	streamName *string
}

func init() {
	register("KinesisStream", ListKinesisStreams)
}

func ListKinesisStreams(sess *session.Session) ([]Resource, error) {
	svc := kinesis.New(sess)
	resources := []Resource{}
	var lastStreamName *string
	params := &kinesis.ListStreamsInput{
		Limit: aws.Int64(25),
	}

	for {
		output, err := svc.ListStreams(params)
		if err != nil {
			return nil, err
		}

		for _, streamName := range output.StreamNames {
			resources = append(resources, &KinesisStream{
				svc:        svc,
				streamName: streamName,
			})
			lastStreamName = streamName
		}

		if *output.HasMoreStreams == false {
			break
		}

		params.ExclusiveStartStreamName = lastStreamName
	}

	return resources, nil
}

func (f *KinesisStream) Remove() error {

	_, err := f.svc.DeleteStream(&kinesis.DeleteStreamInput{
		StreamName: f.streamName,
	})

	return err
}

func (f *KinesisStream) String() string {
	return *f.streamName
}
