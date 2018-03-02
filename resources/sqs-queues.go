package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSQueue struct {
	svc      *sqs.SQS
	queueURL *string
}

func init() {
	register("SQSQueue", ListSQSQueues)
}

func ListSQSQueues(sess *session.Session) ([]Resource, error) {
	svc := sqs.New(sess)

	params := &sqs.ListQueuesInput{}
	resp, err := svc.ListQueues(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, queue := range resp.QueueUrls {
		resources = append(resources, &SQSQueue{
			svc:      svc,
			queueURL: queue,
		})
	}

	return resources, nil
}

func (f *SQSQueue) Remove() error {

	_, err := f.svc.DeleteQueue(&sqs.DeleteQueueInput{
		QueueUrl: f.queueURL,
	})

	return err
}

func (f *SQSQueue) String() string {
	return *f.queueURL
}
