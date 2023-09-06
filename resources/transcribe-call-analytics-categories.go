package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type TranscribeCallAnalyticsCategory struct {
	svc            *transcribeservice.TranscribeService
	name           *string
	inputType      *string
	createTime     *time.Time
	lastUpdateTime *time.Time
}

func init() {
	register("TranscribeCallAnalyticsCategory", ListTranscribeCallAnalyticsCategories)
}

func ListTranscribeCallAnalyticsCategories(sess *session.Session) ([]Resource, error) {
	svc := transcribeservice.New(sess)
	resources := []Resource{}
	var nextToken *string

	for {
		listCallAnalyticsCategoriesInput := &transcribeservice.ListCallAnalyticsCategoriesInput{
			NextToken: nextToken,
		}

		listOutput, err := svc.ListCallAnalyticsCategories(listCallAnalyticsCategoriesInput)
		if err != nil {
			return nil, err
		}
		for _, category := range listOutput.Categories {
			resources = append(resources, &TranscribeCallAnalyticsCategory{
				svc:            svc,
				name:           category.CategoryName,
				inputType:      category.InputType,
				createTime:     category.CreateTime,
				lastUpdateTime: category.LastUpdateTime,
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

func (category *TranscribeCallAnalyticsCategory) Remove() error {
	deleteInput := &transcribeservice.DeleteCallAnalyticsCategoryInput{
		CategoryName: category.name,
	}
	_, err := category.svc.DeleteCallAnalyticsCategory(deleteInput)
	return err
}

func (category *TranscribeCallAnalyticsCategory) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", category.name)
	properties.Set("InputType", category.inputType)
	if category.createTime != nil {
		properties.Set("CreateTime", category.createTime.Format(time.RFC3339))
	}
	if category.lastUpdateTime != nil {
		properties.Set("LastUpdateTime", category.lastUpdateTime.Format(time.RFC3339))
	}
	return properties
}
