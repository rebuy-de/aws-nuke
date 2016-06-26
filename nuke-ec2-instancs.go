package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2Instance struct {
	svc   *ec2.EC2
	id    *string
	state string
}

func (n *EC2Nuke) ListInstances() ([]Resource, error) {
	params := &ec2.DescribeInstancesInput{}
	resp, err := n.svc.DescribeInstances(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			resources = append(resources, &EC2Instance{
				svc:   n.svc,
				id:    instance.InstanceId,
				state: *instance.State.Name,
			})
		}
	}

	return resources, nil
}

func (i *EC2Instance) Check() error {
	if i.state == "terminated" {
		return fmt.Errorf("already terminated")
	}
	return nil
}

func (i *EC2Instance) Remove() error {
	params := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{i.id},
	}

	_, err := i.svc.TerminateInstances(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *EC2Instance) Wait() error {
	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{i.id},
	}
	return i.svc.WaitUntilInstanceTerminated(params)
}

func (i *EC2Instance) String() string {
	return *i.id
}
