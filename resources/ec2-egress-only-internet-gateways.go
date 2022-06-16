package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2EgressOnlyInternetGateway struct {
	svc *ec2.EC2
	igw *ec2.EgressOnlyInternetGateway
}

func init() {
	register("EC2EgressOnlyInternetGateway", ListEC2EgressOnlyInternetGateways)
}

func ListEC2EgressOnlyInternetGateways(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	resources := make([]Resource, 0)
	igwInputParams := &ec2.DescribeEgressOnlyInternetGatewaysInput{
		MaxResults: aws.Int64(255),
	}

	for {
		resp, err := svc.DescribeEgressOnlyInternetGateways(igwInputParams)
		if err != nil {
			return nil, err
		}

		for _, igw := range resp.EgressOnlyInternetGateways {
			resources = append(resources, &EC2EgressOnlyInternetGateway{
				svc: svc,
				igw: igw,
			})
		}

		if resp.NextToken == nil {
			break
		}

		igwInputParams.NextToken = resp.NextToken
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
