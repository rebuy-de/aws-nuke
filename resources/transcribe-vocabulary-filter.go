package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type TranscribeVocabularyFilter struct {
	svc              *transcribeservice.TranscribeService
	name             *string
	languageCode     *string
	lastModifiedTime *time.Time
}

func init() {
	register("TranscribeVocabularyFilter", ListTranscribeVocabularyFilters)
}

func ListTranscribeVocabularyFilters(sess *session.Session) ([]Resource, error) {
	svc := transcribeservice.New(sess)
	resources := []Resource{}
	var nextToken *string

	for {
		listVocabularyFiltersInput := &transcribeservice.ListVocabularyFiltersInput{
			MaxResults: aws.Int64(100),
			NextToken:  nextToken,
		}

		listOutput, err := svc.ListVocabularyFilters(listVocabularyFiltersInput)
		if err != nil {
			return nil, err
		}
		for _, filter := range listOutput.VocabularyFilters {
			resources = append(resources, &TranscribeVocabularyFilter{
				svc:              svc,
				name:             filter.VocabularyFilterName,
				languageCode:     filter.LanguageCode,
				lastModifiedTime: filter.LastModifiedTime,
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

func (filter *TranscribeVocabularyFilter) Remove() error {
	deleteInput := &transcribeservice.DeleteVocabularyFilterInput{
		VocabularyFilterName: filter.name,
	}
	_, err := filter.svc.DeleteVocabularyFilter(deleteInput)
	return err
}

func (filter *TranscribeVocabularyFilter) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", filter.name)
	properties.Set("LanguageCode", filter.languageCode)
	if filter.lastModifiedTime != nil {
		properties.Set("LastModifiedTime", filter.lastModifiedTime.Format(time.RFC3339))
	}
	return properties
}
