package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2InternetGateway struct {
	svc        *ec2.EC2
	igw        *ec2.InternetGateway
	defaultVPC bool
}

func init() {
	register("EC2InternetGateway", ListEC2InternetGateways)
}

func ListEC2InternetGateways(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeInternetGateways(nil)
	if err != nil {
		return nil, err
	}

	defVpcId := ""
	if defVpc := DefaultVpc(svc); defVpc != nil {
		defVpcId = *defVpc.VpcId
	}

	resources := make([]Resource, 0)
	for _, igw := range resp.InternetGateways {
		resources = append(resources, &EC2InternetGateway{
			svc:        svc,
			igw:        igw,
			defaultVPC: HasVpcAttachment(&defVpcId, igw.Attachments),
		})
	}

	return resources, nil
}

func HasVpcAttachment(vpcId *string, attachments []*ec2.InternetGatewayAttachment) bool {
	if *vpcId == "" {
		return false
	}

	for _, attach := range attachments {
		if *vpcId == *attach.VpcId {
			return true
		}
	}
	return false
}

func (e *EC2InternetGateway) Remove() error {
	params := &ec2.DeleteInternetGatewayInput{
		InternetGatewayId: e.igw.InternetGatewayId,
	}

	_, err := e.svc.DeleteInternetGateway(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2InternetGateway) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range e.igw.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	properties.Set("DefaultVPC", e.defaultVPC)
	properties.Set("OwnerID", e.igw.OwnerId)
	return properties
}

func (e *EC2InternetGateway) String() string {
	return *e.igw.InternetGatewayId
}
