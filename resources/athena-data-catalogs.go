package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AthenaDataCatalog struct {
	svc  *athena.Athena
	name *string
}

func init() {
	register("AthenaDataCatalog", ListAthenaDataCatalogs)
}

func ListAthenaDataCatalogs(sess *session.Session) ([]Resource, error) {
	svc := athena.New(sess)
	resources := []Resource{}

	params := &athena.ListDataCatalogsInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.ListDataCatalogs(params)
		if err != nil {
			return nil, err
		}

		for _, catalog := range output.DataCatalogsSummary {
			resources = append(resources, &AthenaDataCatalog{
				svc:  svc,
				name: catalog.CatalogName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *AthenaDataCatalog) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", f.name)

	return properties
}

func (f *AthenaDataCatalog) Remove() error {

	_, err := f.svc.DeleteDataCatalog(&athena.DeleteDataCatalogInput{
		Name: f.name,
	})

	return err
}

func (f *AthenaDataCatalog) Filter() error {
	if *f.name == "AwsDataCatalog" {
		return fmt.Errorf("cannot delete default data source")
	}
	return nil
}

func (f *AthenaDataCatalog) String() string {
	return *f.name
}
