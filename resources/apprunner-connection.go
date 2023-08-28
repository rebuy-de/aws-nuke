package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apprunner"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AppRunnerConnection struct {
	svc            *apprunner.AppRunner
	ConnectionArn  *string
	ConnectionName *string
}

func init() {
	register("AppRunnerConnection", ListAppRunnerConnections)
}

func ListAppRunnerConnections(sess *session.Session) ([]Resource, error) {
	svc := apprunner.New(sess)
	resources := []Resource{}

	params := &apprunner.ListConnectionsInput{}

	for {
		resp, err := svc.ListConnections(params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.ConnectionSummaryList {
			resources = append(resources, &AppRunnerConnection{
				svc:            svc,
				ConnectionArn:  item.ConnectionArn,
				ConnectionName: item.ConnectionName,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *AppRunnerConnection) Remove() error {
	_, err := f.svc.DeleteConnection(&apprunner.DeleteConnectionInput{
		ConnectionArn: f.ConnectionArn,
	})

	return err
}

func (f *AppRunnerConnection) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ConnectionArn", f.ConnectionArn)
	properties.Set("ConnectionName", f.ConnectionName)
	return properties
}
