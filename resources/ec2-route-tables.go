package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2RouteTable struct {
	svc *ec2.EC2
	id  *string
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
			svc: svc,
			id:  out.RouteTableId,
		})
	}

	return resources, nil
}

func (e *EC2RouteTable) Remove() error {
	params := &ec2.DeleteRouteTableInput{
		RouteTableId: e.id,
	}

	_, err := e.svc.DeleteRouteTable(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2RouteTable) String() string {
	return *e.id
}
