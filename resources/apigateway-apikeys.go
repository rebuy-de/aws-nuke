package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
)

type APIGatewayAPIKey struct {
	svc    *apigateway.APIGateway
	APIKey *string
}

func init() {
	register("APIGatewayAPIKey", ListAPIGatewayAPIKeys,
		mapCloudControl("AWS::ApiGateway::ApiKey"))
}

func ListAPIGatewayAPIKeys(sess *session.Session) ([]Resource, error) {
	svc := apigateway.New(sess)
	resources := []Resource{}

	params := &apigateway.GetApiKeysInput{
		Limit: aws.Int64(100),
	}

	for {
		output, err := svc.GetApiKeys(params)
		if err != nil {
			return nil, err
		}

		for _, item := range output.Items {
			resources = append(resources, &APIGatewayAPIKey{
				svc:    svc,
				APIKey: item.Id,
			})
		}

		if output.Position == nil {
			break
		}

		params.Position = output.Position
	}

	return resources, nil
}

func (f *APIGatewayAPIKey) Remove() error {

	_, err := f.svc.DeleteApiKey(&apigateway.DeleteApiKeyInput{
		ApiKey: f.APIKey,
	})

	return err
}

func (f *APIGatewayAPIKey) String() string {
	return *f.APIKey
}
