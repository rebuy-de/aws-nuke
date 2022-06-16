package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func init() {
	register("ComprehendEntitiesDetectionJob", ListComprehendEntitiesDetectionJobs)
}

func ListComprehendEntitiesDetectionJobs(sess *session.Session) ([]Resource, error) {
	svc := comprehend.New(sess)

	params := &comprehend.ListEntitiesDetectionJobsInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListEntitiesDetectionJobs(params)
		if err != nil {
			return nil, err
		}
		for _, entitiesDetectionJob := range resp.EntitiesDetectionJobPropertiesList {
			if *entitiesDetectionJob.JobStatus == "STOPPED" ||
				*entitiesDetectionJob.JobStatus == "FAILED" {
				// if the job has already been stopped, do not try to delete it again
				continue
			}
			resources = append(resources, &ComprehendEntitiesDetectionJob{
				svc:                  svc,
				entitiesDetectionJob: entitiesDetectionJob,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type ComprehendEntitiesDetectionJob struct {
	svc                  *comprehend.Comprehend
	entitiesDetectionJob *comprehend.EntitiesDetectionJobProperties
}

func (ce *ComprehendEntitiesDetectionJob) Remove() error {
	_, err := ce.svc.StopEntitiesDetectionJob(&comprehend.StopEntitiesDetectionJobInput{
		JobId: ce.entitiesDetectionJob.JobId,
	})
	return err
}

func (ce *ComprehendEntitiesDetectionJob) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("JobName", ce.entitiesDetectionJob.JobName)
	properties.Set("JobId", ce.entitiesDetectionJob.JobId)

	return properties
}

func (ce *ComprehendEntitiesDetectionJob) String() string {
	return *ce.entitiesDetectionJob.JobName
}
