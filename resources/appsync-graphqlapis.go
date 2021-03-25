package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appsync"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

// AppSyncGraphqlApi - An AWS AppSync GraphQL API
type AppSyncGraphqlApi struct {
	svc   *appsync.AppSync
	apiID *string
	name  *string
	tags  map[string]*string
}

func init() {
	register("AppSyncGraphqlApi", ListAppSyncGraphqlApis)
}

// ListAppSyncGraphqlApis - List all AWS AppSync GraphQL APIs in the account
func ListAppSyncGraphqlApis(sess *session.Session) ([]Resource, error) {
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

		for _, graphqlApi := range resp.GraphqlApis {
			resources = append(resources, &AppSyncGraphqlApi{
				svc:   svc,
				apiID: graphqlApi.ApiId,
				name:  graphqlApi.Name,
				tags:  graphqlApi.Tags,
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
func (f *AppSyncGraphqlApi) Remove() error {
	_, err := f.svc.DeleteGraphqlApi(&appsync.DeleteGraphqlApiInput{
		ApiId: f.apiID,
	})
	return err
}

// Properties - Get the properties of an AWS AppSync GraphQL API
func (f *AppSyncGraphqlApi) Properties() types.Properties {
	properties := types.NewProperties()
	for key, value := range f.tags {
		properties.SetTag(aws.String(key), value)
	}
	properties.Set("Name", f.name)
	return properties
}

func (f *AppSyncGraphqlApi) String() string {
	return *f.apiID
}
