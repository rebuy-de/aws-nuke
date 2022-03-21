package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type APIGatewayRestAPI struct {
	svc       *apigateway.APIGateway
	restAPIID *string
	name      *string
	version   *string
	tags      map[string]*string
}

func init() {
	register("APIGatewayRestAPI", ListAPIGatewayRestApis)
}

func ListAPIGatewayRestApis(sess *session.Session) ([]Resource, error) {
	svc := apigateway.New(sess)
	resources := []Resource{}

	params := &apigateway.GetRestApisInput{
		Limit: aws.Int64(100),
	}

	for {
		output, err := svc.GetRestApis(params)
		if err != nil {
			return nil, err
		}

		for _, item := range output.Items {
			resources = append(resources, &APIGatewayRestAPI{
				svc:       svc,
				restAPIID: item.Id,
				name:      item.Name,
				version:   item.Version,
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

func (f *APIGatewayRestAPI) Remove() error {

	_, err := f.svc.DeleteRestApi(&apigateway.DeleteRestApiInput{
		RestApiId: f.restAPIID,
	})

	return err
}

func (f *APIGatewayRestAPI) String() string {
	return *f.restAPIID
}

func (f *APIGatewayRestAPI) Properties() types.Properties {
	properties := types.NewProperties()
	for key, tag := range f.tags {
		properties.SetTag(&key, tag)
	}
	properties.
		Set("APIID", f.restAPIID).
		Set("Name", f.name).
		Set("Version", f.version)
	return properties
}
