package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ServiceCatalogProvisionedProduct struct {
	svc            *servicecatalog.ServiceCatalog
	ID             *string
	terminateToken *string
	name           *string
	productID      *string
}

func init() {
	register("ServiceCatalogProvisionedProduct", ListServiceCatalogProvisionedProducts)
}

func ListServiceCatalogProvisionedProducts(sess *session.Session) ([]Resource, error) {
	svc := servicecatalog.New(sess)
	resources := []Resource{}

	params := &servicecatalog.ScanProvisionedProductsInput{
		PageSize: aws.Int64(20),
		AccessLevelFilter: &servicecatalog.AccessLevelFilter{
			Key:   aws.String("Account"),
			Value: aws.String("self"),
		},
	}

	for {
		resp, err := svc.ScanProvisionedProducts(params)
		if err != nil {
			return nil, err
		}

		for _, provisionedProduct := range resp.ProvisionedProducts {
			resources = append(resources, &ServiceCatalogProvisionedProduct{
				svc:            svc,
				ID:             provisionedProduct.Id,
				terminateToken: provisionedProduct.IdempotencyToken,
				name:           provisionedProduct.Name,
				productID:      provisionedProduct.ProductId,
			})
		}

		if resp.NextPageToken == nil {
			break
		}

		params.PageToken = resp.NextPageToken
	}

	return resources, nil
}

func (f *ServiceCatalogProvisionedProduct) Remove() error {

	_, err := f.svc.TerminateProvisionedProduct(&servicecatalog.TerminateProvisionedProductInput{
		ProvisionedProductId: f.ID,
		IgnoreErrors:         aws.Bool(true),
		TerminateToken:       f.terminateToken,
	})

	return err
}

func (f *ServiceCatalogProvisionedProduct) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ID", f.ID)
	properties.Set("Name", f.name)
	properties.Set("ProductID", f.productID)
	return properties
}

func (f *ServiceCatalogProvisionedProduct) String() string {
	return *f.ID
}
