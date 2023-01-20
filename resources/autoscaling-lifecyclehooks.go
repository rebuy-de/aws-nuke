package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func init() {
	register("LifecycleHook", ListLifecycleHooks,
		mapCloudControl("AWS::AutoScaling::LifecycleHook"))
}

func ListLifecycleHooks(s *session.Session) ([]Resource, error) {
	svc := autoscaling.New(s)

	asgResp, err := svc.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{})
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, asg := range asgResp.AutoScalingGroups {
		lchResp, err := svc.DescribeLifecycleHooks(&autoscaling.DescribeLifecycleHooksInput{
			AutoScalingGroupName: asg.AutoScalingGroupName,
		})
		if err != nil {
			return nil, err
		}

		for _, lch := range lchResp.LifecycleHooks {
			resources = append(resources, &LifecycleHook{
				svc:                  svc,
				lifecycleHookName:    lch.LifecycleHookName,
				autoScalingGroupName: lch.AutoScalingGroupName,
			})
		}
	}

	return resources, nil
}

type LifecycleHook struct {
	svc                  *autoscaling.AutoScaling
	lifecycleHookName    *string
	autoScalingGroupName *string
}

func (lch *LifecycleHook) Remove() error {
	params := &autoscaling.DeleteLifecycleHookInput{
		AutoScalingGroupName: lch.autoScalingGroupName,
		LifecycleHookName:    lch.lifecycleHookName,
	}

	_, err := lch.svc.DeleteLifecycleHook(params)
	if err != nil {
		return err
	}

	return nil
}

func (lch *LifecycleHook) String() string {
	return *lch.lifecycleHookName
}
