package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

type AppConfigConfigurationProfile struct {
	svc           *appconfig.AppConfig
	applicationId *string
	id            *string
	name          *string
}

func init() {
	register("AppConfigConfigurationProfile", ListAppConfigConfigurationProfiles)
}

func ListAppConfigConfigurationProfiles(sess *session.Session) ([]Resource, error) {
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
		params := &appconfig.ListConfigurationProfilesInput{
			ApplicationId: application.id,
			MaxResults:    aws.Int64(100),
		}
		err := svc.ListConfigurationProfilesPages(params, func(page *appconfig.ListConfigurationProfilesOutput, lastPage bool) bool {
			for _, item := range page.Items {
				resources = append(resources, &AppConfigConfigurationProfile{
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

func (f *AppConfigConfigurationProfile) Remove() error {
	_, err := f.svc.DeleteConfigurationProfile(&appconfig.DeleteConfigurationProfileInput{
		ApplicationId:          f.applicationId,
		ConfigurationProfileId: f.id,
	})
	return err
}

func (f *AppConfigConfigurationProfile) Properties() types.Properties {
	return types.NewProperties().
		Set("ApplicationID", f.applicationId).
		Set("ID", f.id).
		Set("Name", f.name)
}
