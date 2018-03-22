package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediaconvert"
)

type MediaConvertQueue struct {
	svc  *mediaconvert.MediaConvert
	name *string
}

func init() {
	register("MediaConvertQueue", ListMediaConvertQueues)
}

func ListMediaConvertQueues(sess *session.Session) ([]Resource, error) {
	svc := mediaconvert.New(sess)
	resources := []Resource{}
	var mediaEndpoint *string

	output, err := svc.DescribeEndpoints(&mediaconvert.DescribeEndpointsInput{})
	if err != nil {
		return nil, err
	}

	for _, mediaconvert := range output.Endpoints {
		mediaEndpoint = mediaconvert.Url
	}

	// Update svc to use custom media endpoint
	svc.Endpoint = *mediaEndpoint

	params := &mediaconvert.ListQueuesInput{
		MaxResults: aws.Int64(20),
	}

	for {
		output, err := svc.ListQueues(params)
		if err != nil {
			return nil, err
		}

		for _, queue := range output.Queues {
			resources = append(resources, &MediaConvertQueue{
				svc:  svc,
				name: queue.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *MediaConvertQueue) Remove() error {

	_, err := f.svc.DeleteQueue(&mediaconvert.DeleteQueueInput{
		Name: f.name,
	})

	return err
}

func (f *MediaConvertQueue) String() string {
	return *f.name
}

func (f *MediaConvertQueue) Filter() error {
	if strings.Contains(*f.name, "Default") {
		return fmt.Errorf("cannot delete default queue")
	}
	return nil
}
