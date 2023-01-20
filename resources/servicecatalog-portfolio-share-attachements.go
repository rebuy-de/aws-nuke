package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ServiceCatalogPortfolioShareAttachment struct {
	svc           *servicecatalog.ServiceCatalog
	portfolioID   *string
	accountID     *string
	portfolioName *string
}

func init() {
	register("ServiceCatalogPortfolioShareAttachment", ListServiceCatalogPortfolioShareAttachments)
}

func ListServiceCatalogPortfolioShareAttachments(sess *session.Session) ([]Resource, error) {
	svc := servicecatalog.New(sess)
	resources := []Resource{}
	portfolios := []*servicecatalog.PortfolioDetail{}

	params := &servicecatalog.ListPortfoliosInput{
		PageSize: aws.Int64(20),
	}

	//List all Portfolios
	for {
		resp, err := svc.ListPortfolios(params)
		if err != nil {
			return nil, err
		}

		for _, portfolioDetail := range resp.PortfolioDetails {
			portfolios = append(portfolios, portfolioDetail)
		}

		if resp.NextPageToken == nil {
			break
		}

		params.PageToken = resp.NextPageToken
	}

	accessParams := &servicecatalog.ListPortfolioAccessInput{}

	// Get all accounts which have shared access to the portfolio
	for _, portfolio := range portfolios {

		accessParams.PortfolioId = portfolio.Id

		resp, err := svc.ListPortfolioAccess(accessParams)
		if err != nil {
			return nil, err
		}

		for _, account := range resp.AccountIds {
			resources = append(resources, &ServiceCatalogPortfolioShareAttachment{
				svc:           svc,
				portfolioID:   portfolio.Id,
				accountID:     account,
				portfolioName: portfolio.DisplayName,
			})
		}

	}

	return resources, nil
}

func (f *ServiceCatalogPortfolioShareAttachment) Remove() error {

	_, err := f.svc.DeletePortfolioShare(&servicecatalog.DeletePortfolioShareInput{
		AccountId:   f.accountID,
		PortfolioId: f.portfolioID,
	})

	return err
}

func (f *ServiceCatalogPortfolioShareAttachment) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("PortfolioID", f.portfolioID)
	properties.Set("PortfolioName", f.portfolioName)
	properties.Set("AccountID", f.accountID)
	return properties
}

func (f *ServiceCatalogPortfolioShareAttachment) String() string {
	return fmt.Sprintf("%s -> %s", *f.portfolioID, *f.accountID)
}
