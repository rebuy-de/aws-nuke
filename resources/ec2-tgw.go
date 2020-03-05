package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type EC2TGW struct {
	svc *ec2.EC2
	tgw *ec2.TransitGateway
}

func init() {
	register("EC2TGW", ListEC2TGWs)
}

func ListEC2TGWs(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	params := &ec2.DescribeTransitGatewaysInput{}
	resources := make([]Resource, 0)
	for {
		resp, err := svc.DescribeTransitGateways(params)
		if err != nil {
			return nil, err
		}

		for _, tgw := range resp.TransitGateways {
			resources = append(resources, &EC2TGW{
				svc: svc,
				tgw: tgw,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params = &ec2.DescribeTransitGatewaysInput{
			NextToken: resp.NextToken,
		}
	}

	return resources, nil
}

func (e *EC2TGW) Remove() error {
	params := &ec2.DeleteTransitGatewayInput{
		TransitGatewayId: e.tgw.TransitGatewayId,
	}

	_, err := e.svc.DeleteTransitGateway(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2TGW) Filter() error {
	if *e.tgw.State == "deleted" {
		return fmt.Errorf("already deleted")
	}

	return nil
}

func (e *EC2TGW) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range e.tgw.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	properties.
		Set("ID", e.tgw.TransitGatewayId).
		Set("OwnerId", e.tgw.OwnerId)

	return properties
}

func (e *EC2TGW) String() string {
	return *e.tgw.TransitGatewayId
}
