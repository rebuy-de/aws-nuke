package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
)

type BatchJobQueue struct {
	svc      *batch.Batch
	jobQueue *string
}

func init() {
	register("BatchJobQueue", ListBatchJobQueues)
}

func ListBatchJobQueues(sess *session.Session) ([]Resource, error) {
	svc := batch.New(sess)
	resources := []Resource{}

	params := &batch.DescribeJobQueuesInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeJobQueues(params)
		if err != nil {
			return nil, err
		}

		for _, queue := range output.JobQueues {
			resources = append(resources, &BatchJobQueue{
				svc:      svc,
				jobQueue: queue.JobQueueName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *BatchJobQueue) Remove() error {

	_, err := f.svc.DeleteJobQueue(&batch.DeleteJobQueueInput{
		JobQueue: f.jobQueue,
	})

	return err
}

func (f *BatchJobQueue) String() string {
	return *f.jobQueue
}
