package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type TranscribeMedicalVocabulary struct {
	svc              *transcribeservice.TranscribeService
	name             *string
	state            *string
	languageCode     *string
	lastModifiedTime *time.Time
}

func init() {
	register("TranscribeMedicalVocabulary", ListTranscribeMedicalVocabularies)
}

func ListTranscribeMedicalVocabularies(sess *session.Session) ([]Resource, error) {
	svc := transcribeservice.New(sess)
	resources := []Resource{}
	var nextToken *string

	for {
		listMedicalVocabulariesInput := &transcribeservice.ListMedicalVocabulariesInput{
			MaxResults: aws.Int64(100),
			NextToken:  nextToken,
		}

		listOutput, err := svc.ListMedicalVocabularies(listMedicalVocabulariesInput)
		if err != nil {
			return nil, err
		}
		for _, vocab := range listOutput.Vocabularies {
			resources = append(resources, &TranscribeMedicalVocabulary{
				svc:              svc,
				name:             vocab.VocabularyName,
				state:            vocab.VocabularyState,
				languageCode:     vocab.LanguageCode,
				lastModifiedTime: vocab.LastModifiedTime,
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

func (vocab *TranscribeMedicalVocabulary) Remove() error {
	deleteInput := &transcribeservice.DeleteMedicalVocabularyInput{
		VocabularyName: vocab.name,
	}
	_, err := vocab.svc.DeleteMedicalVocabulary(deleteInput)
	return err
}

func (vocab *TranscribeMedicalVocabulary) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", vocab.name)
	properties.Set("State", vocab.state)
	properties.Set("LanguageCode", vocab.languageCode)
	if vocab.lastModifiedTime != nil {
		properties.Set("LastModifiedTime", vocab.lastModifiedTime.Format(time.RFC3339))
	}
	return properties
}
