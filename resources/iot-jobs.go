package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type IoTJob struct {
	svc    *iot.IoT
	ID     *string
	status *string
}

func init() {
	register("IoTJob", ListIoTJobs)
}

func ListIoTJobs(sess *session.Session) ([]Resource, error) {
	svc := iot.New(sess)
	resources := []Resource{}

	params := &iot.ListJobsInput{
		MaxResults: aws.Int64(100),
		Status:     aws.String("IN_PROGRESS"),
	}
	for {
		output, err := svc.ListJobs(params)
		if err != nil {
			return nil, err
		}

		for _, job := range output.Jobs {
			resources = append(resources, &IoTJob{
				svc:    svc,
				ID:     job.JobId,
				status: job.Status,
			})
		}
		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *IoTJob) Remove() error {

	_, err := f.svc.CancelJob(&iot.CancelJobInput{
		JobId: f.ID,
	})

	return err
}

func (f *IoTJob) String() string {
	return *f.ID
}
