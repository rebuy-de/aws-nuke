package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type APIGatewayVpcLink struct {
	svc       *apigateway.APIGateway
	vpcLinkID *string
	name      *string
	tags      map[string]*string
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
				name:      item.Name,
				tags:      item.Tags,
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

func (f *APIGatewayVpcLink) Properties() types.Properties {
	properties := types.NewProperties()
	for key, tag := range f.tags {
		properties.SetTag(&key, tag)
	}
	properties.
		Set("VPCLinkID", f.vpcLinkID).
		Set("Name", f.name)
	return properties
}
