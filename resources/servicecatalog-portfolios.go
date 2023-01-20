package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ServiceCatalogPortfolio struct {
	svc          *servicecatalog.ServiceCatalog
	ID           *string
	displayName  *string
	providerName *string
}

func init() {
	register("ServiceCatalogPortfolio", ListServiceCatalogPortfolios)
}

func ListServiceCatalogPortfolios(sess *session.Session) ([]Resource, error) {
	svc := servicecatalog.New(sess)
	resources := []Resource{}

	params := &servicecatalog.ListPortfoliosInput{
		PageSize: aws.Int64(20),
	}

	for {
		resp, err := svc.ListPortfolios(params)
		if err != nil {
			return nil, err
		}

		for _, portfolioDetail := range resp.PortfolioDetails {
			resources = append(resources, &ServiceCatalogPortfolio{
				svc:          svc,
				ID:           portfolioDetail.Id,
				displayName:  portfolioDetail.DisplayName,
				providerName: portfolioDetail.ProviderName,
			})
		}

		if resp.NextPageToken == nil {
			break
		}

		params.PageToken = resp.NextPageToken
	}

	return resources, nil
}

func (f *ServiceCatalogPortfolio) Remove() error {

	_, err := f.svc.DeletePortfolio(&servicecatalog.DeletePortfolioInput{
		Id: f.ID,
	})

	return err
}

func (f *ServiceCatalogPortfolio) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ID", f.ID)
	properties.Set("DisplayName", f.displayName)
	properties.Set("ProviderName", f.providerName)
	return properties
}

func (f *ServiceCatalogPortfolio) String() string {
	return *f.ID
}
