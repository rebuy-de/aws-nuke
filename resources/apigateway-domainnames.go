package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
)

type APIGatewayDomainName struct {
	svc        *apigateway.APIGateway
	domainName *string
}

func init() {
	register("APIGatewayDomainName", ListAPIGatewayDomainNames)
}

func ListAPIGatewayDomainNames(sess *session.Session) ([]Resource, error) {
	svc := apigateway.New(sess)
	resources := []Resource{}

	params := &apigateway.GetDomainNamesInput{
		Limit: aws.Int64(100),
	}

	for {
		output, err := svc.GetDomainNames(params)
		if err != nil {
			return nil, err
		}

		for _, item := range output.Items {
			resources = append(resources, &APIGatewayDomainName{
				svc:        svc,
				domainName: item.DomainName,
			})
		}

		if output.Position == nil {
			break
		}

		params.Position = output.Position
	}

	return resources, nil
}

func (f *APIGatewayDomainName) Remove() error {

	_, err := f.svc.DeleteDomainName(&apigateway.DeleteDomainNameInput{
		DomainName: f.domainName,
	})

	return err
}

func (f *APIGatewayDomainName) String() string {
	return *f.domainName
}
