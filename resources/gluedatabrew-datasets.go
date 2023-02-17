package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/gluedatabrew"
)

type GlueDataBrewDatasets struct {
	svc  *gluedatabrew.GlueDataBrew
	name *string
}

func init() {
	register("GlueDataBrewDatasets", ListGlueDatasets)
}

func ListGlueDatasets(sess *session.Session) ([]Resource, error) {
	svc := gluedatabrew.New(sess)
	resources := []Resource{}

	params := &gluedatabrew.ListDatasetsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListDatasets(params)
		if err != nil {
			return nil, err
		}

		for _, dataset := range output.Datasets {
			resources = append(resources, &GlueDataBrewDatasets{
				svc:  svc,
				name: dataset.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueDataBrewDatasets) Remove() error {
	_, err := f.svc.DeleteDataset(&gluedatabrew.DeleteDatasetInput{
		Name: f.name,
	})

	return err
}

func (f *GlueDataBrewDatasets) String() string {
	return *f.name
}
