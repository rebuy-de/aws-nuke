package resources

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type EC2RouteTable struct {
	svc        ec2iface.EC2API
	routeTable *ec2.RouteTable
}

func init() {
	register("EC2RouteTable", ListEC2RouteTables)
}

func ListEC2RouteTables(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeRouteTables(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.RouteTables {
		resources = append(resources, &EC2RouteTable{
			svc:        svc,
			routeTable: out,
		})
	}

	return resources, nil
}

func (e *EC2RouteTable) Remove() error {
	params := &ec2.DeleteRouteTableInput{
		RouteTableId: e.routeTable.RouteTableId,
	}

	_, err := e.svc.DeleteRouteTable(params)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "InvalidRouteTableID.NotFound" { //the route table was deleted elsewhere
				return nil
			}
		}
		return err
	}

	return nil
}

func (e *EC2RouteTable) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range e.routeTable.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	return properties
}

func (e *EC2RouteTable) String() string {
	return *e.routeTable.RouteTableId
}
