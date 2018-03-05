package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
)

type BatchJobQueueState struct {
	svc      *batch.Batch
	jobQueue *string
	state    *string
}

func init() {
	register("BatchJobQueueState", ListBatchJobQueueStates)
}

func ListBatchJobQueueStates(sess *session.Session) ([]Resource, error) {
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
			resources = append(resources, &BatchJobQueueState{
				svc:      svc,
				jobQueue: queue.JobQueueName,
				state:    queue.State,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *BatchJobQueueState) Remove() error {

	_, err := f.svc.UpdateJobQueue(&batch.UpdateJobQueueInput{
		JobQueue: f.jobQueue,
		State:    aws.String("DISABLED"),
	})

	return err
}

func (f *BatchJobQueueState) String() string {
	return *f.jobQueue
}

func (f *BatchJobQueueState) Filter() error {
	if strings.ToLower(*f.state) == "disabled" {
		return fmt.Errorf("already disabled")
	}
	return nil
}
