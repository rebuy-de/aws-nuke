package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	log "github.com/sirupsen/logrus"
)

type ServiceCatalogConstraintPortfolioAttachment struct {
	svc           *servicecatalog.ServiceCatalog
	constraintID  *string
	portfolioID   *string
	portfolioName *string
}

func init() {
	register("ServiceCatalogConstraintPortfolioAttachment", ListServiceCatalogPrincipalProductAttachments)
}

func ListServiceCatalogPrincipalProductAttachments(sess *session.Session) ([]Resource, error) {
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
			if IsAWSError(err, servicecatalog.ErrCodeTagOptionNotMigratedException) {
				log.Info(err)
				break
			}
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

	constraintParams := &servicecatalog.ListConstraintsForPortfolioInput{
		PageSize: aws.Int64(20),
	}

	for _, portfolio := range portfolios {

		constraintParams.PortfolioId = portfolio.Id
		resp, err := svc.ListConstraintsForPortfolio(constraintParams)
		if err != nil {
			return nil, err
		}

		for _, constraintDetail := range resp.ConstraintDetails {
			resources = append(resources, &ServiceCatalogConstraintPortfolioAttachment{
				svc:           svc,
				portfolioID:   portfolio.Id,
				constraintID:  constraintDetail.ConstraintId,
				portfolioName: portfolio.DisplayName,
			})
		}

		if resp.NextPageToken == nil {
			break
		}

		constraintParams.PageToken = resp.NextPageToken
	}

	return resources, nil
}

func (f *ServiceCatalogConstraintPortfolioAttachment) Remove() error {

	_, err := f.svc.DeleteConstraint(&servicecatalog.DeleteConstraintInput{
		Id: f.constraintID,
	})

	return err
}

func (f *ServiceCatalogConstraintPortfolioAttachment) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("PortfolioID", f.portfolioID)
	properties.Set("ConstraintID", f.constraintID)
	properties.Set("PortfolioName", f.portfolioName)
	return properties
}

func (f *ServiceCatalogConstraintPortfolioAttachment) String() string {
	return fmt.Sprintf("%s -> %s", *f.constraintID, *f.portfolioID)
}
