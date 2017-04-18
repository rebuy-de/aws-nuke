package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2RouteTable struct {
	svc    *ec2.EC2
	id     *string
	region *string
}

func (n *EC2Nuke) ListRouteTables() ([]Resource, error) {
	resp, err := n.Service.DescribeRouteTables(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.RouteTables {
		resources = append(resources, &EC2RouteTable{
			svc:    n.Service,
			id:     out.RouteTableId,
			region: n.Service.Config.Region,
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
	return fmt.Sprintf("%s in %s", *e.id, *e.region)
}
