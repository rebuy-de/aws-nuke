package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func init() {
	register("ComprehendTargetedSentimentDetectionJob", ListComprehendTargetedSentimentDetectionJobs)
}

func ListComprehendTargetedSentimentDetectionJobs(sess *session.Session) ([]Resource, error) {
	svc := comprehend.New(sess)

	params := &comprehend.ListTargetedSentimentDetectionJobsInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListTargetedSentimentDetectionJobs(params)
		if err != nil {
			return nil, err
		}
		for _, targetedSentimentDetectionJob := range resp.TargetedSentimentDetectionJobPropertiesList {
			switch *targetedSentimentDetectionJob.JobStatus {
			case "STOPPED", "FAILED", "COMPLETED":
				// if the job has already been stopped, failed, or completed; do not try to stop it again
				continue
			}
			resources = append(resources, &ComprehendTargetedSentimentDetectionJob{
				svc:                           svc,
				targetedSentimentDetectionJob: targetedSentimentDetectionJob,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type ComprehendTargetedSentimentDetectionJob struct {
	svc                           *comprehend.Comprehend
	targetedSentimentDetectionJob *comprehend.TargetedSentimentDetectionJobProperties
}

func (ce *ComprehendTargetedSentimentDetectionJob) Remove() error {
	_, err := ce.svc.StopTargetedSentimentDetectionJob(&comprehend.StopTargetedSentimentDetectionJobInput{
		JobId: ce.targetedSentimentDetectionJob.JobId,
	})
	return err
}

func (ce *ComprehendTargetedSentimentDetectionJob) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("JobName", ce.targetedSentimentDetectionJob.JobName)
	properties.Set("JobId", ce.targetedSentimentDetectionJob.JobId)

	return properties
}

func (ce *ComprehendTargetedSentimentDetectionJob) String() string {
	if ce.targetedSentimentDetectionJob.JobName == nil {
		return "Unnamed job"
	} else {
		return *ce.targetedSentimentDetectionJob.JobName
	}
}
