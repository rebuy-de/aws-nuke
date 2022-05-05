package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RDSDBParameterGroup struct {
	svc  *rds.RDS
	name *string
	tags []*rds.Tag
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
		tags, err := svc.ListTagsForResource(&rds.ListTagsForResourceInput{
			ResourceName: parametergroup.DBParameterGroupArn,
		})

		if err != nil {
			continue
		}

		resources = append(resources, &RDSDBParameterGroup{
			svc:  svc,
			name: parametergroup.DBParameterGroupName,
			tags: tags.TagList,
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

func (i *RDSDBParameterGroup) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", i.name)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
