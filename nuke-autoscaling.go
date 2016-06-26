package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type AutoScalingNuke struct {
	svc *autoscaling.AutoScaling
}

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

func (g AutoScalingGroup) Remove() error {
	params := &autoscaling.DeleteAutoScalingGroupInput{
		AutoScalingGroupName: g.name,
		ForceDelete:          aws.Bool(true),
	}

	_, err := g.svc.DeleteAutoScalingGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (g AutoScalingGroup) Wait() error {
	params := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{g.name},
	}
	return g.svc.WaitUntilGroupNotExists(params)
}

func (g AutoScalingGroup) String() string {
	return *g.name
}
