package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func init() {
	register("LaunchConfiguration", ListLaunchConfigurations)
}

func ListLaunchConfigurations(s *session.Session) ([]Resource, error) {
	resources := make([]Resource, 0)
	svc := autoscaling.New(s)

	params := &autoscaling.DescribeLaunchConfigurationsInput{}
	err := svc.DescribeLaunchConfigurationsPages(params,
		func(page *autoscaling.DescribeLaunchConfigurationsOutput, lastPage bool) bool {
			for _, launchconfig := range page.LaunchConfigurations {
				resources = append(resources, &LaunchConfiguration{
					svc:  svc,
					name: launchconfig.LaunchConfigurationName,
				})
			}
			return !lastPage
		})

	if err != nil {
		return nil, err
	}

	return resources, nil
}

type LaunchConfiguration struct {
	svc  *autoscaling.AutoScaling
	name *string
}

func (launchconfiguration *LaunchConfiguration) Remove() error {
	params := &autoscaling.DeleteLaunchConfigurationInput{
		LaunchConfigurationName: launchconfiguration.name,
	}

	_, err := launchconfiguration.svc.DeleteLaunchConfiguration(params)
	if err != nil {
		return err
	}

	return nil
}

func (launchconfiguration *LaunchConfiguration) String() string {
	return *launchconfiguration.name
}
