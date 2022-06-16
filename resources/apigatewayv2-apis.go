package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type APIGatewayV2API struct {
	svc          *apigatewayv2.ApiGatewayV2
	v2APIID      *string
	name         *string
	protocolType *string
	version      *string
	tags         map[string]*string
}

func init() {
	register("APIGatewayV2API", ListAPIGatewayV2APIs)
}

func ListAPIGatewayV2APIs(sess *session.Session) ([]Resource, error) {
	svc := apigatewayv2.New(sess)
	resources := []Resource{}

	params := &apigatewayv2.GetApisInput{
		MaxResults: aws.String("100"),
	}

	for {
		output, err := svc.GetApis(params)
		if err != nil {
			return nil, err
		}

		for _, item := range output.Items {
			resources = append(resources, &APIGatewayV2API{
				svc:          svc,
				v2APIID:      item.ApiId,
				name:         item.Name,
				protocolType: item.ProtocolType,
				version:      item.Version,
				tags:         item.Tags,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *APIGatewayV2API) Remove() error {

	_, err := f.svc.DeleteApi(&apigatewayv2.DeleteApiInput{
		ApiId: f.v2APIID,
	})

	return err
}

func (f *APIGatewayV2API) String() string {
	return *f.v2APIID
}

func (f *APIGatewayV2API) Properties() types.Properties {
	properties := types.NewProperties()
	for key, tag := range f.tags {
		properties.SetTag(&key, tag)
	}
	properties.
		Set("APIID", f.v2APIID).
		Set("Name", f.name).
		Set("ProtocolType", f.protocolType).
		Set("Version", f.version)
	return properties
}
