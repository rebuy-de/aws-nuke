package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	log "github.com/sirupsen/logrus"
)

type EC2RouteTable struct {
	svc        *ec2.EC2
	routeTable *ec2.RouteTable
	defaultVPC bool
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

	defVpcId := ""
	if defVpc := DefaultVpc(svc); defVpc != nil {
		defVpcId = *defVpc.VpcId
	}

	resources := make([]Resource, 0)
	for _, out := range resp.RouteTables {

		if IsMain(out) {
			log.WithFields(log.Fields{
				"routetable": *out.RouteTableId,
			}).Debug("Main RouteTables cannot be deleted, skipping RouteTable: routetable")
			continue
		}

		resources = append(resources, &EC2RouteTable{
			svc:        svc,
			routeTable: out,
			defaultVPC: defVpcId == *out.VpcId,
		})
	}

	return resources, nil
}

func IsMain(e *ec2.RouteTable) bool {
	for _, association := range e.Associations {
		if *association.Main {
			return true
		}
	}
	return false
}

func (e *EC2RouteTable) Remove() error {
	params := &ec2.DeleteRouteTableInput{
		RouteTableId: e.routeTable.RouteTableId,
	}

	_, err := e.svc.DeleteRouteTable(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2RouteTable) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range e.routeTable.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	properties.Set("DefaultVPC", e.defaultVPC)
	return properties
}

func (e *EC2RouteTable) String() string {
	return *e.routeTable.RouteTableId
}
