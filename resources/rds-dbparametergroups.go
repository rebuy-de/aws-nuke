package resources

import (
	"fmt"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RDSDBParameterGroup struct {
	svc  *rds.RDS
	name *string
}

func (n *RDSNuke) ListParameterGroups() ([]Resource, error) {
	params := &rds.DescribeDBParameterGroupsInput{MaxRecords: aws.Int64(100)}
	resp, err := n.Service.DescribeDBParameterGroups(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, parametergroup := range resp.DBParameterGroups {
		resources = append(resources, &RDSDBParameterGroup{
			svc:  n.Service,
			name: parametergroup.DBParameterGroupName,
		})

	}

	return resources, nil
}

func (i *RDSDBParameterGroup) Filter() error {
	if strings.HasPrefix(*i.name, "default.") {
		return fmt.Errorf("Cannot delete default parameter group")
	}
	return nil
}

func (i *RDSDBParameterGroup) Remove() error {
	params := &rds.DeleteDBParameterGroupInput{
		DBParameterGroupName: i.name,
	}

	_, err := i.svc.DeleteDBParameterGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *RDSDBParameterGroup) String() string {
	return *i.name
}
