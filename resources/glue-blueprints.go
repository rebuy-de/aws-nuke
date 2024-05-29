package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type GlueBlueprint struct {
	svc  *glue.Glue
	name *string
}

func init() {
	register("GlueBlueprint", ListGlueBlueprints)
}

func ListGlueBlueprints(sess *session.Session) ([]Resource, error) {
	svc := glue.New(sess)
	resources := []Resource{}

	params := &glue.ListBlueprintsInput{
		MaxResults: aws.Int64(25),
	}

	for {
		output, err := svc.ListBlueprints(params)
		if err != nil {
			return nil, err
		}

		for _, blueprint := range output.Blueprints {
			resources = append(resources, &GlueBlueprint{
				svc:  svc,
				name: blueprint,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueBlueprint) Remove() error {
	_, err := f.svc.DeleteBlueprint(&glue.DeleteBlueprintInput{
		Name: f.name,
	})

	return err
}

func (f *GlueBlueprint) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", f.name)

	return properties
}

func (f *GlueBlueprint) String() string {
	return *f.name
}
