package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RDSDBSubnetGroup struct {
	svc  *rds.RDS
	name *string
}

func init() {
	register("RDSDBSubnetGroup", ListRDSSubnetGroups)
}

func ListRDSSubnetGroups(sess *session.Session) ([]Resource, error) {
	svc := rds.New(sess)

	params := &rds.DescribeDBSubnetGroupsInput{MaxRecords: aws.Int64(100)}
	resp, err := svc.DescribeDBSubnetGroups(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, subnetGroup := range resp.DBSubnetGroups {
		resources = append(resources, &RDSDBSubnetGroup{
			svc:  svc,
			name: subnetGroup.DBSubnetGroupName,
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
	return *i.name
}
