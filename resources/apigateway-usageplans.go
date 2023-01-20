package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type APIGatewayUsagePlan struct {
	svc         *apigateway.APIGateway
	usagePlanID *string
	name        *string
	tags        map[string]*string
}

func init() {
	register("APIGatewayUsagePlan", ListAPIGatewayUsagePlans,
		mapCloudControl("AWS::ApiGateway::UsagePlan"))
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
				name:        item.Name,
				tags:        item.Tags,
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

func (f *APIGatewayUsagePlan) Properties() types.Properties {
	properties := types.NewProperties()

	for key, tag := range f.tags {
		properties.SetTag(&key, tag)
	}

	properties.
		Set("UsagePlanID", f.usagePlanID).
		Set("Name", f.name)
	return properties
}
