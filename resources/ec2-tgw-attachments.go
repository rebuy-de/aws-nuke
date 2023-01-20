package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2TGWAttachment struct {
	svc  *ec2.EC2
	tgwa *ec2.TransitGatewayAttachment
}

func init() {
	register("EC2TGWAttachment", ListEC2TGWAttachments)
}

func ListEC2TGWAttachments(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	params := &ec2.DescribeTransitGatewayAttachmentsInput{}
	resources := make([]Resource, 0)
	for {
		resp, err := svc.DescribeTransitGatewayAttachments(params)
		if err != nil {
			return nil, err
		}

		for _, tgwa := range resp.TransitGatewayAttachments {
			resources = append(resources, &EC2TGWAttachment{
				svc:  svc,
				tgwa: tgwa,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params = &ec2.DescribeTransitGatewayAttachmentsInput{
			NextToken: resp.NextToken,
		}
	}

	return resources, nil
}

func (e *EC2TGWAttachment) Remove() error {
	if *e.tgwa.ResourceType == "VPN" {
		// This will get deleted as part of EC2VPNConnection, there is no API
		// as part of TGW to delete VPN attachments.
		return fmt.Errorf("VPN attachment")
	}
	params := &ec2.DeleteTransitGatewayVpcAttachmentInput{
		TransitGatewayAttachmentId: e.tgwa.TransitGatewayAttachmentId,
	}

	_, err := e.svc.DeleteTransitGatewayVpcAttachment(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2TGWAttachment) Filter() error {
	if *e.tgwa.State == "deleted" {
		return fmt.Errorf("already deleted")
	}

	return nil
}

func (e *EC2TGWAttachment) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range e.tgwa.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	properties.Set("ID", e.tgwa.TransitGatewayAttachmentId)
	return properties
}

func (e *EC2TGWAttachment) String() string {
	return fmt.Sprintf("%s(%s)", *e.tgwa.TransitGatewayAttachmentId, *e.tgwa.ResourceType)
}
