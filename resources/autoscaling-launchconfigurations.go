package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func init() {
	register("LaunchConfiguration", ListLaunchConfigurations)
}

func ListLaunchConfigurations(s *session.Session) ([]Resource, error) {
	svc := autoscaling.New(s)

	params := &autoscaling.DescribeLaunchConfigurationsInput{}
	resp, err := svc.DescribeLaunchConfigurations(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, launchconfig := range resp.LaunchConfigurations {
		resources = append(resources, &LaunchConfiguration{
			svc:  svc,
			name: launchconfig.LaunchConfigurationName,
		})
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
