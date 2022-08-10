package resources

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2SpotInstanceRequest struct {
	svc                 *ec2.EC2
	spotInstanceRequest *ec2.SpotInstanceRequest
	id                  string
	state               string
}

func init() {
	register("EC2SpotInstanceRequest", ListEC2SpotInstanceRequests)
}

func ListEC2SpotInstanceRequests(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeSpotInstanceRequests(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, config := range resp.SpotInstanceRequests {
		resources = append(resources, &EC2SpotInstanceRequest{
			svc: svc,
			id:  *config.SpotInstanceRequestId,
		})
	}

	return resources, nil
}

func (i *EC2SpotInstanceRequest) Filter() error {
	if i.state == "cancelled" {
		return fmt.Errorf("already cancelled")
	}
	return nil
}

func (i *EC2SpotInstanceRequest) Remove() error {
	params := &ec2.CancelSpotInstanceRequestsInput{
		// TerminateInstances: aws.Bool(true),
		SpotInstanceRequestIds: []*string{
			&i.id,
		},
	}

	_, err := i.svc.CancelSpotInstanceRequests(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *EC2SpotInstanceRequest) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("AvailabilityZoneGroup", i.spotInstanceRequest.AvailabilityZoneGroup)
	properties.Set("InstanceId", i.spotInstanceRequest.InstanceId)
	properties.Set("SpotInstanceRequestId", i.spotInstanceRequest.SpotInstanceRequestId)
	properties.Set("SpotPrice", i.spotInstanceRequest.SpotPrice)
	properties.Set("SpotPrice", i.spotInstanceRequest.SpotPrice)
	properties.Set("SpotInstanceState", i.spotInstanceRequest.State)
	properties.Set("LaunchTime", i.spotInstanceRequest.CreateTime.Format(time.RFC1123))
	properties.Set("SpotInstanceType", i.spotInstanceRequest.Type)

	for _, tagValue := range i.spotInstanceRequest.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}

	return properties
}

func (i *EC2SpotInstanceRequest) String() string {
	return i.id
}
