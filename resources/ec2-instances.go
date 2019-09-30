package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type EC2Instance struct {
	svc      *ec2.EC2
	instance *ec2.Instance
}

func init() {
	register("EC2Instance", ListEC2Instances)
}

func ListEC2Instances(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	params := &ec2.DescribeInstancesInput{}
	resources := make([]Resource, 0)
	for {
		resp, err := svc.DescribeInstances(params)
		if err != nil {
			return nil, err
		}

		for _, reservation := range resp.Reservations {
			for _, instance := range reservation.Instances {
				resources = append(resources, &EC2Instance{
					svc:      svc,
					instance: instance,
				})
			}
		}

		if resp.NextToken == nil {
			break
		}

		params = &ec2.DescribeInstancesInput{
			NextToken: resp.NextToken,
		}
	}

	return resources, nil
}

func (i *EC2Instance) Filter() error {
	if *i.instance.State.Name == "terminated" {
		return fmt.Errorf("already terminated")
	}
	return nil
}

func (i *EC2Instance) Remove() error {
	params := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{i.instance.InstanceId},
	}

	_, err := i.svc.TerminateInstances(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *EC2Instance) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range i.instance.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	return properties
}

func (i *EC2Instance) String() string {
	return *i.instance.InstanceId
}
