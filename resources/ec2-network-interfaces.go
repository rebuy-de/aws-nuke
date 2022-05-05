package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2NetworkInterface struct {
	svc *ec2.EC2
	eni *ec2.NetworkInterface
}

func init() {
	register("EC2NetworkInterface", ListEC2NetworkInterfaces)
}

func ListEC2NetworkInterfaces(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeNetworkInterfaces(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.NetworkInterfaces {

		resources = append(resources, &EC2NetworkInterface{
			svc: svc,
			eni: out,
		})
	}

	return resources, nil
}

func (e *EC2NetworkInterface) Remove() error {

	if e.eni.Attachment != nil {
		_, err := e.svc.DetachNetworkInterface(&ec2.DetachNetworkInterfaceInput{
			AttachmentId: e.eni.Attachment.AttachmentId,
			Force:        aws.Bool(true),
		})
		if err != nil {
			if e.eni.Attachment.AttachmentId != nil {
				expected := fmt.Sprintf("The interface attachment '%s' does not exist.", *e.eni.Attachment.AttachmentId)
				if !strings.Contains(err.Error(), expected) {
					return err
				}
			}

		}
	}

	params := &ec2.DeleteNetworkInterfaceInput{
		NetworkInterfaceId: e.eni.NetworkInterfaceId,
	}

	_, err := e.svc.DeleteNetworkInterface(params)
	if err != nil {
		return err
	}

	return nil
}

func (r *EC2NetworkInterface) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tag := range r.eni.TagSet {
		properties.SetTag(tag.Key, tag.Value)
	}
	properties.
		Set("ID", r.eni.NetworkInterfaceId).
		Set("VPC", r.eni.VpcId).
		Set("AvailabilityZone", r.eni.AvailabilityZone).
		Set("PrivateIPAddress", r.eni.PrivateIpAddress).
		Set("SubnetID", r.eni.SubnetId).
		Set("Status", r.eni.Status)
	return properties
}

func (r *EC2NetworkInterface) String() string {
	return *r.eni.NetworkInterfaceId
}
