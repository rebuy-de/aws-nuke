package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type TranscribeLanguageModel struct {
	svc                 *transcribeservice.TranscribeService
	name                *string
	baseModelName       *string
	createTime          *time.Time
	failureReason       *string
	languageCode        *string
	lastModifiedTime    *time.Time
	modelStatus         *string
	upgradeAvailability *bool
}

func init() {
	register("TranscribeLanguageModel", ListTranscribeLanguageModels)
}

func ListTranscribeLanguageModels(sess *session.Session) ([]Resource, error) {
	svc := transcribeservice.New(sess)
	resources := []Resource{}
	var nextToken *string

	for {
		listLanguageModelsInput := &transcribeservice.ListLanguageModelsInput{
			MaxResults: aws.Int64(100),
			NextToken:  nextToken,
		}

		listOutput, err := svc.ListLanguageModels(listLanguageModelsInput)
		if err != nil {
			return nil, err
		}
		for _, model := range listOutput.Models {
			resources = append(resources, &TranscribeLanguageModel{
				svc:                 svc,
				name:                model.ModelName,
				baseModelName:       model.BaseModelName,
				createTime:          model.CreateTime,
				failureReason:       model.FailureReason,
				languageCode:        model.LanguageCode,
				lastModifiedTime:    model.LastModifiedTime,
				modelStatus:         model.ModelStatus,
				upgradeAvailability: model.UpgradeAvailability,
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

func (model *TranscribeLanguageModel) Remove() error {
	deleteInput := &transcribeservice.DeleteLanguageModelInput{
		ModelName: model.name,
	}
	_, err := model.svc.DeleteLanguageModel(deleteInput)
	return err
}

func (model *TranscribeLanguageModel) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", model.name)
	properties.Set("BaseModelName", model.baseModelName)
	if model.createTime != nil {
		properties.Set("CreateTime", model.createTime.Format(time.RFC3339))
	}
	properties.Set("FailureReason", model.failureReason)
	properties.Set("LanguageCode", model.languageCode)
	if model.lastModifiedTime != nil {
		properties.Set("LastModifiedTime", model.lastModifiedTime.Format(time.RFC3339))
	}
	properties.Set("ModelStatus", model.modelStatus)
	properties.Set("UpgradeAvailability", model.upgradeAvailability)
	return properties
}
