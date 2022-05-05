package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func init() {
	register("ComprehendDominantLanguageDetectionJob", ListComprehendDominantLanguageDetectionJobs)
}

func ListComprehendDominantLanguageDetectionJobs(sess *session.Session) ([]Resource, error) {
	svc := comprehend.New(sess)

	params := &comprehend.ListDominantLanguageDetectionJobsInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListDominantLanguageDetectionJobs(params)
		if err != nil {
			return nil, err
		}
		for _, dominantLanguageDetectionJob := range resp.DominantLanguageDetectionJobPropertiesList {
			resources = append(resources, &ComprehendDominantLanguageDetectionJob{
				svc:                          svc,
				dominantLanguageDetectionJob: dominantLanguageDetectionJob,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type ComprehendDominantLanguageDetectionJob struct {
	svc                          *comprehend.Comprehend
	dominantLanguageDetectionJob *comprehend.DominantLanguageDetectionJobProperties
}

func (ce *ComprehendDominantLanguageDetectionJob) Remove() error {
	_, err := ce.svc.StopDominantLanguageDetectionJob(&comprehend.StopDominantLanguageDetectionJobInput{
		JobId: ce.dominantLanguageDetectionJob.JobId,
	})
	return err
}

func (ce *ComprehendDominantLanguageDetectionJob) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("JobName", ce.dominantLanguageDetectionJob.JobName)
	properties.Set("JobId", ce.dominantLanguageDetectionJob.JobId)

	return properties
}

func (ce *ComprehendDominantLanguageDetectionJob) String() string {
	return *ce.dominantLanguageDetectionJob.JobName
}
