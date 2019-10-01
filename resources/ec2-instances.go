package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			resources = append(resources, &EC2Instance{
				svc:      svc,
				instance: instance,
			})
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
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "OperationNotPermitted" &&

			awsErr.Message() == "The instance '"+*i.instance.InstanceId+"' may not be terminated. "+
				"Modify its 'disableApiTermination' instance attribute and try again." {
			err = i.DisableProtection()
			if err != nil {
				return err
			}
			_, err := i.svc.TerminateInstances(params)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

func (i *EC2Instance) DisableProtection() error {
	params := &ec2.ModifyInstanceAttributeInput{
		InstanceId: i.instance.InstanceId,
		DisableApiTermination: &ec2.AttributeBooleanValue{
			Value: aws.Bool(false),
		},
	}
	_, err := i.svc.ModifyInstanceAttribute(params)
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
