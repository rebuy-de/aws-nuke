package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appsync"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

// AppSyncGraphqlAPI - An AWS AppSync GraphQL API
type AppSyncGraphqlAPI struct {
	svc   *appsync.AppSync
	apiID *string
	name  *string
	tags  map[string]*string
}

func init() {
	register("AppSyncGraphqlAPI", ListAppSyncGraphqlAPIs)
}

// ListAppSyncGraphqlAPIs - List all AWS AppSync GraphQL APIs in the account
func ListAppSyncGraphqlAPIs(sess *session.Session) ([]Resource, error) {
	svc := appsync.New(sess)
	resources := []Resource{}

	params := &appsync.ListGraphqlApisInput{
		MaxResults: aws.Int64(25),
	}

	for {
		resp, err := svc.ListGraphqlApis(params)
		if err != nil {
			return nil, err
		}

		for _, graphqlAPI := range resp.GraphqlApis {
			resources = append(resources, &AppSyncGraphqlAPI{
				svc:   svc,
				apiID: graphqlAPI.ApiId,
				name:  graphqlAPI.Name,
				tags:  graphqlAPI.Tags,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

// Remove - remove an AWS AppSync GraphQL API
func (f *AppSyncGraphqlAPI) Remove() error {
	_, err := f.svc.DeleteGraphqlApi(&appsync.DeleteGraphqlApiInput{
		ApiId: f.apiID,
	})
	return err
}

// Properties - Get the properties of an AWS AppSync GraphQL API
func (f *AppSyncGraphqlAPI) Properties() types.Properties {
	properties := types.NewProperties()
	for key, value := range f.tags {
		properties.SetTag(aws.String(key), value)
	}
	properties.Set("Name", f.name)
	properties.Set("APIID", f.apiID)
	return properties
}

func (f *AppSyncGraphqlAPI) String() string {
	return *f.apiID
}
