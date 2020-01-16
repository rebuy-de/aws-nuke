package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/robomaker"
)

type RoboMakerSimulationJob struct {
	svc  *robomaker.RoboMaker
	name *string
	arn  *string
}

func init() {
	register("RoboMakerSimulationJob", ListRoboMakerSimulationJobs)
}

func simulationJobNeedsToBeCanceled(job *robomaker.SimulationJobSummary) bool {
	for _, n := range []string{"Completed", "Failed", "RunningFailed", "Terminating", "Terminated", "Canceled"} {
		if job.Status != nil && *job.Status == n {
			return false
		}
	}
	return true
}
func ListRoboMakerSimulationJobs(sess *session.Session) ([]Resource, error) {
	svc := robomaker.New(sess)
	resources := []Resource{}

	params := &robomaker.ListSimulationJobsInput{
		MaxResults: aws.Int64(30),
	}

	for {
		resp, err := svc.ListSimulationJobs(params)
		if err != nil {
			return nil, err
		}

		for _, simulationJob := range resp.SimulationJobSummaries {
			if simulationJobNeedsToBeCanceled(simulationJob) {
				resources = append(resources, &RoboMakerSimulationJob{
					svc: svc,
					arn: simulationJob.Arn,
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

func (f *RoboMakerSimulationJob) Remove() error {

	_, err := f.svc.CancelSimulationJob(&robomaker.CancelSimulationJobInput{
		Job: f.arn,
	})

	return err
}

func (f *RoboMakerSimulationJob) String() string {
	return *f.arn
}
