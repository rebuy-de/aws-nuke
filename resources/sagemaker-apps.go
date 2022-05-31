package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type SageMakerApp struct {
	svc             *sagemaker.SageMaker
	domainID        *string
	appName         *string
	appType         *string
	userProfileName *string
	status          *string
}

func init() {
	register("SageMakerApp", ListSageMakerApps)
}

func ListSageMakerApps(sess *session.Session) ([]Resource, error) {
	svc := sagemaker.New(sess)
	resources := []Resource{}

	params := &sagemaker.ListAppsInput{
		MaxResults: aws.Int64(30),
	}

	for {
		resp, err := svc.ListApps(params)
		if err != nil {
			return nil, err
		}

		for _, app := range resp.Apps {
			resources = append(resources, &SageMakerApp{
				svc:             svc,
				domainID:        app.DomainId,
				appName:         app.AppName,
				appType:         app.AppType,
				userProfileName: app.UserProfileName,
				status:          app.Status,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *SageMakerApp) Remove() error {
	_, err := f.svc.DeleteApp(&sagemaker.DeleteAppInput{
		DomainId:        f.domainID,
		AppName:         f.appName,
		AppType:         f.appType,
		UserProfileName: f.userProfileName,
	})

	return err
}

func (f *SageMakerApp) String() string {
	return *f.appName
}

func (i *SageMakerApp) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("DomainID", i.domainID).
		Set("AppName", i.appName).
		Set("AppType", i.appType).
		Set("UserProfileName", i.userProfileName)
	return properties
}

func (f *SageMakerApp) Filter() error {
	if *f.status == sagemaker.AppStatusDeleted {
		return fmt.Errorf("already deleted")
	}
	return nil
}
