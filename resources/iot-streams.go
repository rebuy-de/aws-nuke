package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTStream struct {
	svc *iot.IoT
	ID  *string
}

func init() {
	register("IoTStream", ListIoTStreams)
}

func ListIoTStreams(sess *session.Session) ([]Resource, error) {
	svc := iot.New(sess)
	resources := []Resource{}

	params := &iot.ListStreamsInput{
		MaxResults: aws.Int64(100),
	}
	for {
		output, err := svc.ListStreams(params)
		if err != nil {
			return nil, err
		}

		for _, stream := range output.Streams {
			resources = append(resources, &IoTStream{
				svc: svc,
				ID:  stream.StreamId,
			})
		}
		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *IoTStream) Remove() error {

	_, err := f.svc.DeleteStream(&iot.DeleteStreamInput{
		StreamId: f.ID,
	})

	return err
}

func (f *IoTStream) String() string {
	return *f.ID
}
