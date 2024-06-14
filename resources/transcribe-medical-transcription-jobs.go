package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type TranscribeMedicalTranscriptionJob struct {
	svc                       *transcribeservice.TranscribeService
	name                      *string
	status                    *string
	completionTime            *time.Time
	contentIdentificationType *string
	creationTime              *time.Time
	failureReason             *string
	languageCode              *string
	outputLocationType        *string
	specialty                 *string
	startTime                 *time.Time
	inputType                 *string
}

func init() {
	register("TranscribeMedicalTranscriptionJob", ListTranscribeMedicalTranscriptionJobs)
}

func ListTranscribeMedicalTranscriptionJobs(sess *session.Session) ([]Resource, error) {
	svc := transcribeservice.New(sess)
	resources := []Resource{}
	var nextToken *string

	for {
		listMedicalTranscriptionJobsInput := &transcribeservice.ListMedicalTranscriptionJobsInput{
			MaxResults: aws.Int64(100),
			NextToken:  nextToken,
		}

		listOutput, err := svc.ListMedicalTranscriptionJobs(listMedicalTranscriptionJobsInput)
		if err != nil {
			return nil, err
		}
		for _, job := range listOutput.MedicalTranscriptionJobSummaries {
			resources = append(resources, &TranscribeMedicalTranscriptionJob{
				svc:                       svc,
				name:                      job.MedicalTranscriptionJobName,
				status:                    job.TranscriptionJobStatus,
				completionTime:            job.CompletionTime,
				contentIdentificationType: job.ContentIdentificationType,
				creationTime:              job.CreationTime,
				failureReason:             job.FailureReason,
				languageCode:              job.LanguageCode,
				outputLocationType:        job.OutputLocationType,
				specialty:                 job.Specialty,
				startTime:                 job.StartTime,
				inputType:                 job.Type,
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

func (job *TranscribeMedicalTranscriptionJob) Remove() error {
	deleteInput := &transcribeservice.DeleteMedicalTranscriptionJobInput{
		MedicalTranscriptionJobName: job.name,
	}
	_, err := job.svc.DeleteMedicalTranscriptionJob(deleteInput)
	return err
}

func (job *TranscribeMedicalTranscriptionJob) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", job.name)
	properties.Set("Status", job.status)
	if job.completionTime != nil {
		properties.Set("CompletionTime", job.completionTime.Format(time.RFC3339))
	}
	properties.Set("ContentIdentificationType", job.contentIdentificationType)
	if job.creationTime != nil {
		properties.Set("CreationTime", job.creationTime.Format(time.RFC3339))
	}
	properties.Set("FailureReason", job.failureReason)
	properties.Set("LanguageCode", job.languageCode)
	properties.Set("OutputLocationType", job.outputLocationType)
	properties.Set("Specialty", job.specialty)
	if job.startTime != nil {
		properties.Set("StartTime", job.startTime.Format(time.RFC3339))
	}
	properties.Set("InputType", job.inputType)
	return properties
}
