package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apprunner"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AppRunnerService struct {
	svc         *apprunner.AppRunner
	ServiceArn  *string
	ServiceId   *string
	ServiceName *string
}

func init() {
	register("AppRunnerService", ListAppRunnerServices)
}

func ListAppRunnerServices(sess *session.Session) ([]Resource, error) {
	svc := apprunner.New(sess)
	resources := []Resource{}

	params := &apprunner.ListServicesInput{}

	for {
		resp, err := svc.ListServices(params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.ServiceSummaryList {
			resources = append(resources, &AppRunnerService{
				svc:         svc,
				ServiceArn:  item.ServiceArn,
				ServiceId:   item.ServiceId,
				ServiceName: item.ServiceName,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *AppRunnerService) Remove() error {
	_, err := f.svc.DeleteService(&apprunner.DeleteServiceInput{
		ServiceArn: f.ServiceArn,
	})

	return err
}

func (f *AppRunnerService) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ServiceArn", f.ServiceArn)
	properties.Set("ServiceId", f.ServiceId)
	properties.Set("ServiceName", f.ServiceName)
	return properties
}
