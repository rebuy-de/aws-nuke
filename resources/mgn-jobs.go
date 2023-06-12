package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mgn"
)

type MgnJob struct {
	svc *mgn.Mgn
	id  *string
}

func init() {
	register("MgnJob", ListMgnJobs)
}

func ListMgnJobs(sess *session.Session) ([]Resource, error) {
	svc := mgn.New(sess)
	resources := []Resource{}

	params := &mgn.DescribeJobsInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.DescribeJobs(params)
		if err != nil {
			return nil, err
		}

		for _, job := range output.Items {
			resources = append(resources, &MgnJob{
				svc: svc,
				id:  job.JobID,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *MgnJob) Remove() error {

	_, err := f.svc.CancelJob(&mgn.CancelJobInput{
		JobID: f.id,
	})

	return err
}

func (f *MgnJob) String() string {
	return *f.id
}
