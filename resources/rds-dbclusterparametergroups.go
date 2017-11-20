package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RDSDBClusterParameterGroup struct {
	svc  *rds.RDS
	name *string
}

func (n *RDSNuke) ListClusterParameterGroups() ([]Resource, error) {
	params := &rds.DescribeDBClusterParameterGroupsInput{MaxRecords: aws.Int64(100)}
	resp, err := n.Service.DescribeDBClusterParameterGroups(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, parametergroup := range resp.DBClusterParameterGroups {
		resources = append(resources, &RDSDBClusterParameterGroup{
			svc:  n.Service,
			name: parametergroup.DBClusterParameterGroupName,
		})

	}

	return resources, nil
}

func (i *RDSDBClusterParameterGroup) Filter() error {
	if strings.HasPrefix(*i.name, "default.") {
		return fmt.Errorf("Cannot delete default parameter group")
	}
	return nil
}

func (i *RDSDBClusterParameterGroup) Remove() error {
	params := &rds.DeleteDBClusterParameterGroupInput{
		DBClusterParameterGroupName: i.name,
	}

	_, err := i.svc.DeleteDBClusterParameterGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *RDSDBClusterParameterGroup) String() string {
	return *i.name
}
