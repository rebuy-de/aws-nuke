package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CodePipelineCustomActionType struct {
	svc      *codepipeline.CodePipeline
	owner    *string
	category *string
	provider *string
}

func init() {
	register("CodePipelineCustomActionType", ListCodePipelineCustomActionTypes)
}

func ListCodePipelineCustomActionTypes(sess *session.Session) ([]Resource, error) {
	svc := codepipeline.New(sess)
	resources := []Resource{}

	params := &codepipeline.ListActionTypesInput{}

	for {
		resp, err := svc.ListActionTypes(params)
		if err != nil {
			return nil, err
		}

		for _, actionTypes := range resp.ActionTypes {
			resources = append(resources, &CodePipelineCustomActionType{
				svc:      svc,
				owner:    actionTypes.Id.Owner,
				category: actionTypes.Id.Category,
				provider: actionTypes.Id.Provider,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *CodePipelineCustomActionType) Filter() error {
	if !strings.HasPrefix(*f.owner, "Custom") {
		return fmt.Errorf("cannot delete default codepipeline custom action type")
	}
	return nil
}

func (f *CodePipelineCustomActionType) Remove() error {
	_, err := f.svc.DeleteCustomActionType(&codepipeline.DeleteCustomActionTypeInput{
		Category: f.category,
	})

	return err
}

func (f *CodePipelineCustomActionType) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Category", f.category)
	properties.Set("Owner", f.owner)
	properties.Set("Provider", f.provider)
	return properties
}

func (f *CodePipelineCustomActionType) String() string {
	return *f.owner
}
