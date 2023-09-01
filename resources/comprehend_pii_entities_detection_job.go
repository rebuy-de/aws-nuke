package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func init() {
	register("ComprehendPiiEntititesDetectionJob", ListComprehendPiiEntitiesDetectionJobs)
}

func ListComprehendPiiEntitiesDetectionJobs(sess *session.Session) ([]Resource, error) {
	svc := comprehend.New(sess)

	params := &comprehend.ListPiiEntitiesDetectionJobsInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListPiiEntitiesDetectionJobs(params)
		if err != nil {
			return nil, err
		}
		for _, piiEntititesDetectionJob := range resp.PiiEntitiesDetectionJobPropertiesList {
			switch *piiEntititesDetectionJob.JobStatus {
			case "STOPPED", "FAILED", "COMPLETED":
				// if the job has already been stopped, failed, or completed; do not try to stop it again
				continue
			}
			resources = append(resources, &ComprehendPiiEntitiesDetectionJob{
				svc:                      svc,
				piiEntititesDetectionJob: piiEntititesDetectionJob,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type ComprehendPiiEntitiesDetectionJob struct {
	svc                      *comprehend.Comprehend
	piiEntititesDetectionJob *comprehend.PiiEntitiesDetectionJobProperties
}

func (ce *ComprehendPiiEntitiesDetectionJob) Remove() error {
	_, err := ce.svc.StopPiiEntitiesDetectionJob(&comprehend.StopPiiEntitiesDetectionJobInput{
		JobId: ce.piiEntititesDetectionJob.JobId,
	})
	return err
}

func (ce *ComprehendPiiEntitiesDetectionJob) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("JobName", ce.piiEntititesDetectionJob.JobName)
	properties.Set("JobId", ce.piiEntititesDetectionJob.JobId)

	return properties
}

func (ce *ComprehendPiiEntitiesDetectionJob) String() string {
	if ce.piiEntititesDetectionJob.JobName == nil {
		return "Unnamed job"
	} else {
		return *ce.piiEntititesDetectionJob.JobName
	}
}
