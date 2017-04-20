package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RDSDBSubnetGroup struct {
	svc    *rds.RDS
	name   *string
	region *string
}

func (n *RDSNuke) ListSubnetGroups() ([]Resource, error) {
	params := &rds.DescribeDBSubnetGroupsInput{MaxRecords: aws.Int64(100)}
	resp, err := n.Service.DescribeDBSubnetGroups(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, subnetGroup := range resp.DBSubnetGroups {
		resources = append(resources, &RDSDBSubnetGroup{
			svc:    n.Service,
			name:   subnetGroup.DBSubnetGroupName,
			region: n.Service.Config.Region,
		})

	}

	return resources, nil
}

func (i *RDSDBSubnetGroup) Remove() error {
	params := &rds.DeleteDBSubnetGroupInput{
		DBSubnetGroupName: i.name,
	}

	_, err := i.svc.DeleteDBSubnetGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *RDSDBSubnetGroup) String() string {
	return fmt.Sprintf("%s in %s", *i.name, *i.region)
}
