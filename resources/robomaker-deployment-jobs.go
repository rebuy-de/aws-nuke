package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/robomaker"
)

type RoboMakerDeploymentJob struct {
	svc  *robomaker.RoboMaker
	name *string
	arn  *string
}

func init() {
	register("RoboMakerDeploymentJob", ListRoboMakerDeploymentJobs)
}

func deploymentJobNeedsToBeCanceled(job *robomaker.DeploymentJob) bool {
	for _, n := range []string{"Completed", "Failed", "RunningFailed", "Terminating", "Terminated", "Canceled"} {
		if job.Status != nil && *job.Status == n {
			return false
		}
	}
	return true
}

func ListRoboMakerDeploymentJobs(sess *session.Session) ([]Resource, error) {
	svc := robomaker.New(sess)
	resources := []Resource{}

	params := &robomaker.ListDeploymentJobsInput{
		MaxResults: aws.Int64(30),
	}

	for {
		resp, err := svc.ListDeploymentJobs(params)
		if err != nil {
			return nil, err
		}

		for _, deploymentJob := range resp.DeploymentJobs {
			if deploymentJobNeedsToBeCanceled(deploymentJob) {
				resources = append(resources, &RoboMakerDeploymentJob{
					svc: svc,
					arn: deploymentJob.Arn,
				})
			}
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *RoboMakerDeploymentJob) Remove() error {

	_, err := f.svc.CancelDeploymentJob(&robomaker.CancelDeploymentJobInput{
		Job: f.arn,
	})

	return err
}

func (f *RoboMakerDeploymentJob) String() string {
	return *f.arn
}
