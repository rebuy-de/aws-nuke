package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CodePipelineWebhook struct {
	svc  *codepipeline.CodePipeline
	name *string
}

func init() {
	register("CodePipelineWebhook", ListCodePipelineWebhooks)
}

func ListCodePipelineWebhooks(sess *session.Session) ([]Resource, error) {
	svc := codepipeline.New(sess)
	resources := []Resource{}

	params := &codepipeline.ListWebhooksInput{}

	for {
		resp, err := svc.ListWebhooks(params)
		if err != nil {
			return nil, err
		}

		for _, webHooks := range resp.Webhooks {
			resources = append(resources, &CodePipelineWebhook{
				svc:  svc,
				name: webHooks.Definition.Name,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *CodePipelineWebhook) Remove() error {
	_, err := f.svc.DeleteWebhook(&codepipeline.DeleteWebhookInput{
		Name: f.name,
	})

	return err
}

func (f *CodePipelineWebhook) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", f.name)
	return properties
}

func (f *CodePipelineWebhook) String() string {
	return *f.name
}
