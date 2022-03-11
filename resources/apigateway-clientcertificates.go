package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
)

type APIGatewayClientCertificate struct {
	svc                 *apigateway.APIGateway
	clientCertificateID *string
}

func init() {
	register("APIGatewayClientCertificate", ListAPIGatewayClientCertificates,
		mapCloudControl("AWS::ApiGateway::ClientCertificate"))
}

func ListAPIGatewayClientCertificates(sess *session.Session) ([]Resource, error) {
	svc := apigateway.New(sess)
	resources := []Resource{}

	params := &apigateway.GetClientCertificatesInput{
		Limit: aws.Int64(100),
	}

	for {
		output, err := svc.GetClientCertificates(params)
		if err != nil {
			return nil, err
		}

		for _, item := range output.Items {
			resources = append(resources, &APIGatewayClientCertificate{
				svc:                 svc,
				clientCertificateID: item.ClientCertificateId,
			})
		}

		if output.Position == nil {
			break
		}

		params.Position = output.Position
	}

	return resources, nil
}

func (f *APIGatewayClientCertificate) Remove() error {

	_, err := f.svc.DeleteClientCertificate(&apigateway.DeleteClientCertificateInput{
		ClientCertificateId: f.clientCertificateID,
	})

	return err
}

func (f *APIGatewayClientCertificate) String() string {
	return *f.clientCertificateID
}
