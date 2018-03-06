package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
)

type GlueJob struct {
	svc     *glue.Glue
	jobName *string
}

func init() {
	register("GlueJob", ListGlueJobs)
}

func ListGlueJobs(sess *session.Session) ([]Resource, error) {
	svc := glue.New(sess)
	resources := []Resource{}

	params := &glue.GetJobsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.GetJobs(params)
		if err != nil {
			return nil, err
		}

		for _, job := range output.Jobs {
			resources = append(resources, &GlueJob{
				svc:     svc,
				jobName: job.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueJob) Remove() error {

	_, err := f.svc.DeleteJob(&glue.DeleteJobInput{
		JobName: f.jobName,
	})

	return err
}

func (f *GlueJob) String() string {
	return *f.jobName
}
