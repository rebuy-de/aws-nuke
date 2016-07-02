package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func (n *AutoScalingNuke) ListGroups() ([]Resource, error) {
	params := &autoscaling.DescribeAutoScalingGroupsInput{}
	resp, err := n.svc.DescribeAutoScalingGroups(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, asg := range resp.AutoScalingGroups {
		resources = append(resources, &AutoScalingGroup{
			svc:  n.svc,
			name: asg.AutoScalingGroupName,
		})
	}
	return resources, nil
}

type AutoScalingGroup struct {
	svc  *autoscaling.AutoScaling
	name *string
}

func (asg *AutoScalingGroup) Remove() error {
	params := &autoscaling.DeleteAutoScalingGroupInput{
		AutoScalingGroupName: asg.name,
		ForceDelete:          aws.Bool(true),
	}

	_, err := asg.svc.DeleteAutoScalingGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (asg *AutoScalingGroup) Wait() error {
	params := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{asg.name},
	}
	return asg.svc.WaitUntilGroupNotExists(params)
}

func (asg *AutoScalingGroup) String() string {
	return *asg.name
}
