package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
)

type FirehoseDeliveryStream struct {
	svc                *firehose.Firehose
	deliveryStreamName *string
	tags               []*firehose.Tag
}

func init() {
	register("FirehoseDeliveryStream", ListFirehoseDeliveryStreams)
}

func ListFirehoseDeliveryStreams(sess *session.Session) ([]Resource, error) {
	svc := firehose.New(sess)
	resources := []Resource{}
	tags := []*firehose.Tag{}
	var lastDeliveryStreamName *string

	params := &firehose.ListDeliveryStreamsInput{
		Limit: aws.Int64(25),
	}

	for {
		output, err := svc.ListDeliveryStreams(params)
		if err != nil {
			return nil, err
		}

		for _, deliveryStreamName := range output.DeliveryStreamNames {
			tagParams := &firehose.ListTagsForDeliveryStreamInput{
				DeliveryStreamName: deliveryStreamName,
				Limit:              aws.Int64(1), // todo put 100
			}

			for {
				tagResp, tagErr := svc.ListTagsForDeliveryStream(tagParams)
				if tagErr != nil {
					return nil, tagErr
				}

				tags = append(tags, tagResp.Tags...)
				if !*tagResp.HasMoreTags {
					break
				}
			}

			resources = append(resources, &FirehoseDeliveryStream{
				svc:                svc,
				deliveryStreamName: deliveryStreamName,
				tags:               tags,
			})

			lastDeliveryStreamName = deliveryStreamName
		}

		if !*output.HasMoreDeliveryStreams {
			break
		}

		params.ExclusiveStartDeliveryStreamName = lastDeliveryStreamName
	}

	return resources, nil
}

func (f *FirehoseDeliveryStream) Remove() error {

	_, err := f.svc.DeleteDeliveryStream(&firehose.DeleteDeliveryStreamInput{
		DeliveryStreamName: f.deliveryStreamName,
	})

	return err
}

func (f *FirehoseDeliveryStream) String() string {
	return *f.deliveryStreamName
}
