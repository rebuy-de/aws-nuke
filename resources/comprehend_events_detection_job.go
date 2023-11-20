package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func init() {
	register("ComprehendEventsDetectionJob", ListComprehendEventsDetectionJobs)
}

func ListComprehendEventsDetectionJobs(sess *session.Session) ([]Resource, error) {
	svc := comprehend.New(sess)

	params := &comprehend.ListEventsDetectionJobsInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListEventsDetectionJobs(params)
		if err != nil {
			return nil, err
		}
		for _, eventsDetectionJob := range resp.EventsDetectionJobPropertiesList {
			switch *eventsDetectionJob.JobStatus {
			case "STOPPED", "FAILED", "COMPLETED":
				// if the job has already been stopped, failed, or completed; do not try to stop it again
				continue
			}
			resources = append(resources, &ComprehendEventsDetectionJob{
				svc:                svc,
				eventsDetectionJob: eventsDetectionJob,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type ComprehendEventsDetectionJob struct {
	svc                *comprehend.Comprehend
	eventsDetectionJob *comprehend.EventsDetectionJobProperties
}

func (ce *ComprehendEventsDetectionJob) Remove() error {
	_, err := ce.svc.StopEventsDetectionJob(&comprehend.StopEventsDetectionJobInput{
		JobId: ce.eventsDetectionJob.JobId,
	})
	return err
}

func (ce *ComprehendEventsDetectionJob) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("JobName", ce.eventsDetectionJob.JobName)
	properties.Set("JobId", ce.eventsDetectionJob.JobId)

	return properties
}

func (ce *ComprehendEventsDetectionJob) String() string {
	if ce.eventsDetectionJob.JobName == nil {
		return "Unnamed job"
	} else {
		return *ce.eventsDetectionJob.JobName
	}
}
