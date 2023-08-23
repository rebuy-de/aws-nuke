package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AppConfigApplication struct {
	svc  *appconfig.AppConfig
	id   *string
	name *string
}

func init() {
	register("AppConfigApplication", ListAppConfigApplications)
}

func ListAppConfigApplications(sess *session.Session) ([]Resource, error) {
	svc := appconfig.New(sess)
	resources := []Resource{}
	params := &appconfig.ListApplicationsInput{
		MaxResults: aws.Int64(50),
	}
	err := svc.ListApplicationsPages(params, func(page *appconfig.ListApplicationsOutput, lastPage bool) bool {
		for _, item := range page.Items {
			resources = append(resources, &AppConfigApplication{
				svc:  svc,
				id:   item.Id,
				name: item.Name,
			})
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func (f *AppConfigApplication) Remove() error {
	_, err := f.svc.DeleteApplication(&appconfig.DeleteApplicationInput{
		ApplicationId: f.id,
	})
	return err
}

func (f *AppConfigApplication) Properties() types.Properties {
	return types.NewProperties().
		Set("ID", f.id).
		Set("Name", f.name)
}
