package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type EC2EgressOnlyInternetGateway struct {
	svc ec2iface.EC2API
	igw *ec2.EgressOnlyInternetGateway
}

func init() {
	register("EC2EgressOnlyInternetGateway", ListEC2EgressOnlyInternetGateway)
}

func ListEC2EgressOnlyInternetGateway(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeEgressOnlyInternetGateways(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, igw := range resp.EgressOnlyInternetGateways {
		resources = append(resources, &EC2EgressOnlyInternetGateway{
			svc: svc,
			igw: igw,
		})
	}

	return resources, nil
}

func (e *EC2EgressOnlyInternetGateway) Remove() error {
	params := &ec2.DeleteEgressOnlyInternetGatewayInput{
		EgressOnlyInternetGatewayId: e.igw.EgressOnlyInternetGatewayId,
	}

	_, err := e.svc.DeleteEgressOnlyInternetGateway(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2EgressOnlyInternetGateway) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range e.igw.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	return properties
}

func (e *EC2EgressOnlyInternetGateway) String() string {
	return *e.igw.EgressOnlyInternetGatewayId
}
