package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AppConfigDeploymentStrategy struct {
	svc  *appconfig.AppConfig
	id   *string
	name *string
}

func init() {
	register("AppConfigDeploymentStrategy", ListAppConfigDeploymentStrategies)
}

func ListAppConfigDeploymentStrategies(sess *session.Session) ([]Resource, error) {
	svc := appconfig.New(sess)
	resources := []Resource{}
	params := &appconfig.ListDeploymentStrategiesInput{
		MaxResults: aws.Int64(50),
	}
	err := svc.ListDeploymentStrategiesPages(params, func(page *appconfig.ListDeploymentStrategiesOutput, lastPage bool) bool {
		for _, item := range page.Items {
			resources = append(resources, &AppConfigDeploymentStrategy{
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

func (f *AppConfigDeploymentStrategy) Filter() error {
	if strings.HasPrefix(*f.name, "AppConfig.") {
		return fmt.Errorf("cannot delete predefined Deployment Strategy")
	}
	return nil
}

func (f *AppConfigDeploymentStrategy) Remove() error {
	_, err := f.svc.DeleteDeploymentStrategy(&appconfig.DeleteDeploymentStrategyInput{
		DeploymentStrategyId: f.id,
	})
	return err
}

func (f *AppConfigDeploymentStrategy) Properties() types.Properties {
	return types.NewProperties().
		Set("ID", f.id).
		Set("Name", f.name)
}
