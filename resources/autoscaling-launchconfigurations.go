package resources

import "github.com/aws/aws-sdk-go/service/autoscaling"

func (n *AutoScalingNuke) ListLaunchConfigurations() ([]Resource, error) {
	params := &autoscaling.DescribeLaunchConfigurationsInput{}
	resp, err := n.Service.DescribeLaunchConfigurations(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, launchconfig := range resp.LaunchConfigurations {
		resources = append(resources, &LaunchConfiguration{
			svc:  n.Service,
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
