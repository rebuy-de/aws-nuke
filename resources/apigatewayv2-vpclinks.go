package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type APIGatewayV2VpcLink struct {
	svc       *apigatewayv2.ApiGatewayV2
	vpcLinkID *string
	name      *string
	tags      map[string]*string
}

func init() {
	register("APIGatewayV2VpcLink", ListAPIGatewayV2VpcLinks)
}

func ListAPIGatewayV2VpcLinks(sess *session.Session) ([]Resource, error) {
	svc := apigatewayv2.New(sess)
	resources := []Resource{}

	params := &apigatewayv2.GetVpcLinksInput{
		MaxResults: aws.String("100"),
	}

	for {
		output, err := svc.GetVpcLinks(params)
		if err != nil {
			return nil, err
		}

		for _, item := range output.Items {
			resources = append(resources, &APIGatewayV2VpcLink{
				svc:       svc,
				vpcLinkID: item.VpcLinkId,
				name:      item.Name,
				tags:      item.Tags,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *APIGatewayV2VpcLink) Remove() error {

	_, err := f.svc.DeleteVpcLink(&apigatewayv2.DeleteVpcLinkInput{
		VpcLinkId: f.vpcLinkID,
	})

	return err
}

func (f *APIGatewayV2VpcLink) String() string {
	return *f.vpcLinkID
}

func (f *APIGatewayV2VpcLink) Properties() types.Properties {
	properties := types.NewProperties()
	for key, tag := range f.tags {
		properties.SetTag(&key, tag)
	}
	properties.
		Set("VPCLinkID", f.vpcLinkID).
		Set("Name", f.name)
	return properties
}
