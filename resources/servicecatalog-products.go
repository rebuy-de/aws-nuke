package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ServiceCatalogProduct struct {
	svc  *servicecatalog.ServiceCatalog
	ID   *string
	name *string
}

func init() {
	register("ServiceCatalogProduct", ListServiceCatalogProducts)
}

func ListServiceCatalogProducts(sess *session.Session) ([]Resource, error) {
	svc := servicecatalog.New(sess)
	resources := []Resource{}

	params := &servicecatalog.SearchProductsAsAdminInput{
		PageSize: aws.Int64(20),
	}

	for {
		resp, err := svc.SearchProductsAsAdmin(params)
		if err != nil {
			return nil, err
		}

		for _, productView := range resp.ProductViewDetails {
			resources = append(resources, &ServiceCatalogProduct{
				svc:  svc,
				ID:   productView.ProductViewSummary.ProductId,
				name: productView.ProductViewSummary.Name,
			})
		}

		if resp.NextPageToken == nil {
			break
		}

		params.PageToken = resp.NextPageToken
	}

	return resources, nil
}

func (f *ServiceCatalogProduct) Remove() error {

	_, err := f.svc.DeleteProduct(&servicecatalog.DeleteProductInput{
		Id: f.ID,
	})

	return err
}

func (f *ServiceCatalogProduct) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ID", f.ID)
	properties.Set("Name", f.name)
	return properties
}

func (f *ServiceCatalogProduct) String() string {
	return *f.ID
}
