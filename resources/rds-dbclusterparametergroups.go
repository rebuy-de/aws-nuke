package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RDSDBClusterParameterGroup struct {
	svc  *rds.RDS
	name *string
	tags []*rds.Tag
}

func init() {
	register("RDSDBClusterParameterGroup", ListRDSClusterParameterGroups)
}

func ListRDSClusterParameterGroups(sess *session.Session) ([]Resource, error) {
	svc := rds.New(sess)

	params := &rds.DescribeDBClusterParameterGroupsInput{MaxRecords: aws.Int64(100)}
	resp, err := svc.DescribeDBClusterParameterGroups(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, parametergroup := range resp.DBClusterParameterGroups {
		tags, err := svc.ListTagsForResource(&rds.ListTagsForResourceInput{
                        ResourceName: parametergroup.DBClusterParameterGroupArn,
                })

                if err != nil {
                        continue
                }

		resources = append(resources, &RDSDBClusterParameterGroup{
			svc:  svc,
			name: parametergroup.DBClusterParameterGroupName,
			tags: tags.TagList,
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

func (i *RDSDBClusterParameterGroup) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", i.name)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
