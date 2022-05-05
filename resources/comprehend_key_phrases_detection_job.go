package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func init() {
	register("ComprehendKeyPhrasesDetectionJob", ListComprehendKeyPhrasesDetectionJobs)
}

func ListComprehendKeyPhrasesDetectionJobs(sess *session.Session) ([]Resource, error) {
	svc := comprehend.New(sess)

	params := &comprehend.ListKeyPhrasesDetectionJobsInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListKeyPhrasesDetectionJobs(params)
		if err != nil {
			return nil, err
		}
		for _, keyPhrasesDetectionJob := range resp.KeyPhrasesDetectionJobPropertiesList {
			resources = append(resources, &ComprehendKeyPhrasesDetectionJob{
				svc:                    svc,
				keyPhrasesDetectionJob: keyPhrasesDetectionJob,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type ComprehendKeyPhrasesDetectionJob struct {
	svc                    *comprehend.Comprehend
	keyPhrasesDetectionJob *comprehend.KeyPhrasesDetectionJobProperties
}

func (ce *ComprehendKeyPhrasesDetectionJob) Remove() error {
	_, err := ce.svc.StopKeyPhrasesDetectionJob(&comprehend.StopKeyPhrasesDetectionJobInput{
		JobId: ce.keyPhrasesDetectionJob.JobId,
	})
	return err
}

func (ce *ComprehendKeyPhrasesDetectionJob) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("JobName", ce.keyPhrasesDetectionJob.JobName)
	properties.Set("JobId", ce.keyPhrasesDetectionJob.JobId)

	return properties
}

func (ce *ComprehendKeyPhrasesDetectionJob) String() string {
	return *ce.keyPhrasesDetectionJob.JobName
}
