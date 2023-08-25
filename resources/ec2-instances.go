package resources

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/config"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2Instance struct {
	svc      *ec2.EC2
	instance *ec2.Instance

	featureFlags config.FeatureFlags
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

func (i *EC2Instance) FeatureFlags(ff config.FeatureFlags) {
	i.featureFlags = ff
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
		// Check for Termination Protection, disable it, and try termination again.
		if ok && awsErr.Code() == "OperationNotPermitted" &&
			awsErr.Message() == "The instance '"+*i.instance.InstanceId+"' may not be "+
				"terminated. Modify its 'disableApiTermination' instance attribute and "+
				"try again." && i.featureFlags.DisableDeletionProtection.EC2Instance {
			termErr := i.DisableTerminationProtection()
			if termErr != nil {
				return termErr
			}
			_, err = i.svc.TerminateInstances(params)
			// If we still get an error, we'll check for type and let the next routine
			// handle it.
			if err != nil {
				awsErr, ok = err.(awserr.Error)
			}
		}

		// Check for Stop Protection, disable it, and try termination again.
		if ok && awsErr.Code() == "OperationNotPermitted" &&
			awsErr.Message() == "The instance '"+*i.instance.InstanceId+"' may not be "+
				"terminated. Modify its 'disableApiStop' instance attribute and try "+
				"again." && i.featureFlags.DisableEC2InstanceStopProtection {
			stopErr := i.DisableStopProtection()
			if stopErr != nil {
				return stopErr
			}
			_, err = i.svc.TerminateInstances(params)
		}

		// If we still have an error at this point, we'll return it.
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *EC2Instance) DisableStopProtection() error {
	params := &ec2.ModifyInstanceAttributeInput{
		InstanceId: i.instance.InstanceId,
		DisableApiStop: &ec2.AttributeBooleanValue{
			Value: aws.Bool(false),
		},
	}
	_, err := i.svc.ModifyInstanceAttribute(params)
	if err != nil {
		return err
	}
	return nil
}

func (i *EC2Instance) DisableTerminationProtection() error {
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
	properties.Set("Identifier", i.instance.InstanceId)
	properties.Set("ImageIdentifier", i.instance.ImageId)
	properties.Set("InstanceState", i.instance.State.Name)
	properties.Set("InstanceType", i.instance.InstanceType)
	properties.Set("LaunchTime", i.instance.LaunchTime.Format(time.RFC3339))

	for _, tagValue := range i.instance.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}

	return properties
}

func (i *EC2Instance) String() string {
	return *i.instance.InstanceId
}
