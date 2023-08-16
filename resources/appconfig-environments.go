package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

type AppConfigEnvironment struct {
	svc           *appconfig.AppConfig
	applicationId *string
	id            *string
	name          *string
}

func init() {
	register("AppConfigEnvironment", ListAppConfigEnvironments)
}

func ListAppConfigEnvironments(sess *session.Session) ([]Resource, error) {
	svc := appconfig.New(sess)
	resources := []Resource{}
	applications, err := ListAppConfigApplications(sess)
	if err != nil {
		return nil, err
	}
	for _, applicationResource := range applications {
		application, ok := applicationResource.(*AppConfigApplication)
		if !ok {
			logrus.Errorf("Unable to cast AppConfigApplication.")
			continue
		}
		params := &appconfig.ListEnvironmentsInput{
			ApplicationId: application.id,
			MaxResults:    aws.Int64(50),
		}
		err := svc.ListEnvironmentsPages(params, func(page *appconfig.ListEnvironmentsOutput, lastPage bool) bool {
			for _, item := range page.Items {
				resources = append(resources, &AppConfigEnvironment{
					svc:           svc,
					applicationId: application.id,
					id:            item.Id,
					name:          item.Name,
				})
			}
			return true
		})
		if err != nil {
			return nil, err
		}
	}
	return resources, nil
}

func (f *AppConfigEnvironment) Remove() error {
	_, err := f.svc.DeleteEnvironment(&appconfig.DeleteEnvironmentInput{
		ApplicationId: f.applicationId,
		EnvironmentId: f.id,
	})
	return err
}

func (f *AppConfigEnvironment) Properties() types.Properties {
	return types.NewProperties().
		Set("ApplicationID", f.applicationId).
		Set("ID", f.id).
		Set("Name", f.name)
}
