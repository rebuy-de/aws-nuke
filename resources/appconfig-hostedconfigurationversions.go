package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

type AppConfigHostedConfigurationVersion struct {
	svc                    *appconfig.AppConfig
	applicationId          *string
	configurationProfileId *string
	versionNumber          *int64
}

func init() {
	register("AppConfigHostedConfigurationVersion", ListAppConfigHostedConfigurationVersions)
}

func ListAppConfigHostedConfigurationVersions(sess *session.Session) ([]Resource, error) {
	svc := appconfig.New(sess)
	resources := []Resource{}
	configurationProfiles, err := ListAppConfigConfigurationProfiles(sess)
	if err != nil {
		return nil, err
	}
	for _, configurationProfileResource := range configurationProfiles {
		configurationProfile, ok := configurationProfileResource.(*AppConfigConfigurationProfile)
		if !ok {
			logrus.Errorf("Unable to cast AppConfigConfigurationProfile.")
			continue
		}
		params := &appconfig.ListHostedConfigurationVersionsInput{
			ApplicationId:          configurationProfile.applicationId,
			ConfigurationProfileId: configurationProfile.id,
			MaxResults:             aws.Int64(100),
		}
		err := svc.ListHostedConfigurationVersionsPages(params, func(page *appconfig.ListHostedConfigurationVersionsOutput, lastPage bool) bool {
			for _, item := range page.Items {
				resources = append(resources, &AppConfigHostedConfigurationVersion{
					svc:                    svc,
					applicationId:          configurationProfile.applicationId,
					configurationProfileId: configurationProfile.id,
					versionNumber:          item.VersionNumber,
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

func (f *AppConfigHostedConfigurationVersion) Remove() error {
	_, err := f.svc.DeleteHostedConfigurationVersion(&appconfig.DeleteHostedConfigurationVersionInput{
		ApplicationId:          f.applicationId,
		ConfigurationProfileId: f.configurationProfileId,
		VersionNumber:          f.versionNumber,
	})
	return err
}

func (f *AppConfigHostedConfigurationVersion) Properties() types.Properties {
	return types.NewProperties().
		Set("ApplicationID", f.applicationId).
		Set("ConfigurationProfileID", f.configurationProfileId).
		Set("VersionNumber", f.versionNumber)
}
