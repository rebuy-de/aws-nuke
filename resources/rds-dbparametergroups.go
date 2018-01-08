package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RDSDBParameterGroup struct {
	svc  *rds.RDS
	name *string
}

func init() {
	register("RDSDBParameterGroup", ListRDSParameterGroups)
}

func ListRDSParameterGroups(sess *session.Session) ([]Resource, error) {
	svc := rds.New(sess)

	params := &rds.DescribeDBParameterGroupsInput{MaxRecords: aws.Int64(100)}
	resp, err := svc.DescribeDBParameterGroups(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, parametergroup := range resp.DBParameterGroups {
		resources = append(resources, &RDSDBParameterGroup{
			svc:  svc,
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
