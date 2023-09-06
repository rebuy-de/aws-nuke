package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type TranscribeTranscriptionJob struct {
	svc            *transcribeservice.TranscribeService
	name           *string
	status         *string
	completionTime *time.Time
	creationTime   *time.Time
	failureReason  *string
	languageCode   *string
	startTime      *time.Time
}

func init() {
	register("TranscribeTranscriptionJob", ListTranscribeTranscriptionJobs)
}

func ListTranscribeTranscriptionJobs(sess *session.Session) ([]Resource, error) {
	svc := transcribeservice.New(sess)
	resources := []Resource{}
	var nextToken *string

	for {
		listTranscriptionJobsInput := &transcribeservice.ListTranscriptionJobsInput{
			MaxResults: aws.Int64(100),
			NextToken:  nextToken,
		}

		listOutput, err := svc.ListTranscriptionJobs(listTranscriptionJobsInput)
		if err != nil {
			return nil, err
		}
		for _, job := range listOutput.TranscriptionJobSummaries {
			resources = append(resources, &TranscribeTranscriptionJob{
				svc:            svc,
				name:           job.TranscriptionJobName,
				status:         job.TranscriptionJobStatus,
				completionTime: job.CompletionTime,
				creationTime:   job.CreationTime,
				failureReason:  job.FailureReason,
				languageCode:   job.LanguageCode,
				startTime:      job.StartTime,
			})
		}

		// Check if there are more results
		if listOutput.NextToken == nil {
			break // No more results, exit the loop
		}

		// Set the nextToken for the next iteration
		nextToken = listOutput.NextToken
	}
	return resources, nil
}

func (job *TranscribeTranscriptionJob) Remove() error {
	deleteInput := &transcribeservice.DeleteTranscriptionJobInput{
		TranscriptionJobName: job.name,
	}
	_, err := job.svc.DeleteTranscriptionJob(deleteInput)
	return err
}

func (job *TranscribeTranscriptionJob) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", job.name)
	properties.Set("Status", job.status)
	if job.completionTime != nil {
		properties.Set("CompletionTime", job.completionTime.Format(time.RFC3339))
	}
	if job.creationTime != nil {
		properties.Set("CreationTime", job.creationTime.Format(time.RFC3339))
	}
	properties.Set("FailureReason", job.failureReason)
	properties.Set("LanguageCode", job.languageCode)
	if job.startTime != nil {
		properties.Set("StartTime", job.startTime.Format(time.RFC3339))
	}
	return properties
}
