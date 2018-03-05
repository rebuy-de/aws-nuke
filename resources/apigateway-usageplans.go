package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
)

type APIGatewayUsagePlan struct {
	svc         *apigateway.APIGateway
	usagePlanID *string
}

func init() {
	register("APIGatewayUsagePlan", ListAPIGatewayUsagePlans)
}

func ListAPIGatewayUsagePlans(sess *session.Session) ([]Resource, error) {
	svc := apigateway.New(sess)
	resources := []Resource{}

	params := &apigateway.GetUsagePlansInput{
		Limit: aws.Int64(100),
	}

	for {
		output, err := svc.GetUsagePlans(params)
		if err != nil {
			return nil, err
		}

		for _, item := range output.Items {
			resources = append(resources, &APIGatewayUsagePlan{
				svc:         svc,
				usagePlanID: item.Id,
			})
		}

		if output.Position == nil {
			break
		}

		params.Position = output.Position
	}

	return resources, nil
}

func (f *APIGatewayUsagePlan) Remove() error {

	_, err := f.svc.DeleteUsagePlan(&apigateway.DeleteUsagePlanInput{
		UsagePlanId: f.usagePlanID,
	})

	return err
}

func (f *APIGatewayUsagePlan) String() string {
	return *f.usagePlanID
}
