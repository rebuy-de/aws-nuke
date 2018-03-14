package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
)

type ServiceCatalogPortfolioProductAttachment struct {
	svc         *servicecatalog.ServiceCatalog
	productID   *string
	portfolioID *string
}

func init() {
	register("ServiceCatalogPortfolioProductAttachment", ListServiceCatalogPortfolioProductAttachments)
}

func ListServiceCatalogPortfolioProductAttachments(sess *session.Session) ([]Resource, error) {
	svc := servicecatalog.New(sess)
	resources := []Resource{}
	products := []*string{}

	params := &servicecatalog.SearchProductsAsAdminInput{
		PageSize: aws.Int64(20),
	}

	//List all Products and then search assigned portfolios
	for {
		resp, err := svc.SearchProductsAsAdmin(params)
		if err != nil {
			return nil, err
		}

		for _, productViewDetail := range resp.ProductViewDetails {
			products = append(products, productViewDetail.ProductViewSummary.ProductId)
		}

		if resp.NextPageToken == nil {
			break
		}

		params.PageToken = resp.NextPageToken
	}

	portfolioParams := &servicecatalog.ListPortfoliosForProductInput{
		PageSize: aws.Int64(20),
	}

	for _, product := range products {

		portfolioParams.ProductId = product

		resp, err := svc.ListPortfoliosForProduct(portfolioParams)
		if err != nil {
			return nil, err
		}

		for _, portfolioDetail := range resp.PortfolioDetails {
			resources = append(resources, &ServiceCatalogPortfolioProductAttachment{
				svc:         svc,
				productID:   product,
				portfolioID: portfolioDetail.Id,
			})
		}

		if resp.NextPageToken == nil {
			continue
		}

		portfolioParams.PageToken = resp.NextPageToken
	}

	return resources, nil
}

func (f *ServiceCatalogPortfolioProductAttachment) Remove() error {

	_, err := f.svc.DisassociateProductFromPortfolio(&servicecatalog.DisassociateProductFromPortfolioInput{
		ProductId:   f.productID,
		PortfolioId: f.portfolioID,
	})

	return err
}

func (f *ServiceCatalogPortfolioProductAttachment) String() string {
	return fmt.Sprintf("%s -> %s", *f.productID, *f.portfolioID)
}
