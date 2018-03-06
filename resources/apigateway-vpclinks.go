package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
)

type APIGatewayVpcLink struct {
	svc       *apigateway.APIGateway
	vpcLinkID *string
}

func init() {
	register("APIGatewayVpcLink", ListAPIGatewayVpcLinks)
}

func ListAPIGatewayVpcLinks(sess *session.Session) ([]Resource, error) {
	svc := apigateway.New(sess)
	resources := []Resource{}

	params := &apigateway.GetVpcLinksInput{
		Limit: aws.Int64(100),
	}

	for {
		output, err := svc.GetVpcLinks(params)
		if err != nil {
			return nil, err
		}

		for _, item := range output.Items {
			resources = append(resources, &APIGatewayVpcLink{
				svc:       svc,
				vpcLinkID: item.Id,
			})
		}

		if output.Position == nil {
			break
		}

		params.Position = output.Position
	}

	return resources, nil
}

func (f *APIGatewayVpcLink) Remove() error {

	_, err := f.svc.DeleteVpcLink(&apigateway.DeleteVpcLinkInput{
		VpcLinkId: f.vpcLinkID,
	})

	return err
}

func (f *APIGatewayVpcLink) String() string {
	return *f.vpcLinkID
}
