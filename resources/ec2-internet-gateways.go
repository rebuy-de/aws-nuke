package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type EC2InternetGateway struct {
	svc *ec2.EC2
	igw *ec2.InternetGateway
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

	resources := make([]Resource, 0)
	for _, igw := range resp.InternetGateways {
		resources = append(resources, &EC2InternetGateway{
			svc: svc,
			igw: igw,
		})
	}

	return resources, nil
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
	return properties
}

func (e *EC2InternetGateway) String() string {
	return *e.igw.InternetGatewayId
}
