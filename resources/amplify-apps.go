package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/amplify"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AmplifyApp struct {
	svc   *amplify.Amplify
	appID *string
	name  *string
	tags  map[string]*string
}

func init() {
	register("AmplifyApp", ListAmplifyApps)
}

func ListAmplifyApps(sess *session.Session) ([]Resource, error) {
	svc := amplify.New(sess)
	resources := []Resource{}

	params := &amplify.ListAppsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListApps(params)
		if err != nil {
			return nil, err
		}

		for _, item := range output.Apps {
			resources = append(resources, &AmplifyApp{
				svc:   svc,
				appID: item.AppId,
				name:  item.Name,
				tags:  item.Tags,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *AmplifyApp) Remove() error {
	_, err := f.svc.DeleteApp(&amplify.DeleteAppInput{
		AppId: f.appID,
	})

	return err
}

func (f *AmplifyApp) String() string {
	return *f.appID
}

func (f *AmplifyApp) Properties() types.Properties {
	properties := types.NewProperties()
	for key, tag := range f.tags {
		properties.SetTag(&key, tag)
	}
	properties.
		Set("AppID", f.appID).
		Set("Name", f.name)
	return properties
}
