package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
)

type ServiceCatalogPrincipalPortfolioAttachment struct {
	svc          *servicecatalog.ServiceCatalog
	portfolioID  *string
	principalARN *string
}

func init() {
	register("ServiceCatalogPrincipalPortfolioAttachment", ListServiceCatalogPrincipalPortfolioAttachments)
}

func ListServiceCatalogPrincipalPortfolioAttachments(sess *session.Session) ([]Resource, error) {
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

	principalParams := &servicecatalog.ListPrincipalsForPortfolioInput{
		PageSize: aws.Int64(20),
	}

	for _, portfolio := range portfolios {

		principalParams.PortfolioId = portfolio.Id

		resp, err := svc.ListPrincipalsForPortfolio(principalParams)
		if err != nil {
			return nil, err
		}

		for _, principal := range resp.Principals {
			resources = append(resources, &ServiceCatalogPrincipalPortfolioAttachment{
				svc:          svc,
				principalARN: principal.PrincipalARN,
				portfolioID:  portfolio.Id,
			})
		}

		if resp.NextPageToken == nil {
			break
		}

		principalParams.PageToken = resp.NextPageToken
	}

	return resources, nil
}

func (f *ServiceCatalogPrincipalPortfolioAttachment) Remove() error {

	_, err := f.svc.DisassociatePrincipalFromPortfolio(&servicecatalog.DisassociatePrincipalFromPortfolioInput{
		PrincipalARN: f.principalARN,
		PortfolioId:  f.portfolioID,
	})

	return err
}

func (f *ServiceCatalogPrincipalPortfolioAttachment) String() string {
	return fmt.Sprintf("%s -> %s", *f.principalARN, *f.portfolioID)
}
