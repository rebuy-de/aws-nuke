package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type GlueWorkflow struct {
	svc  *glue.Glue
	name *string
}

func init() {
	register("GlueWorkflow", ListGlueWorkflows)
}

func ListGlueWorkflows(sess *session.Session) ([]Resource, error) {
	svc := glue.New(sess)
	resources := []Resource{}

	params := &glue.ListWorkflowsInput{
		MaxResults: aws.Int64(25),
	}

	for {
		output, err := svc.ListWorkflows(params)
		if err != nil {
			return nil, err
		}

		for _, workflowName := range output.Workflows {
			resources = append(resources, &GlueWorkflow{
				svc:  svc,
				name: workflowName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueWorkflow) Remove() error {
	_, err := f.svc.DeleteWorkflow(&glue.DeleteWorkflowInput{
		Name: f.name,
	})

	return err
}

func (f *GlueWorkflow) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", f.name)

	return properties
}

func (f *GlueWorkflow) String() string {
	return *f.name
}
