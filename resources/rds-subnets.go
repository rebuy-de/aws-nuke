package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RDSDBSubnetGroup struct {
	svc  *rds.RDS
	name *string
	tags []*rds.Tag
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
		tags, err := svc.ListTagsForResource(&rds.ListTagsForResourceInput{
                        ResourceName: subnetGroup.DBSubnetGroupArn,
                })

                if err != nil {
                        continue
                }

		resources = append(resources, &RDSDBSubnetGroup{
			svc:  svc,
			name: subnetGroup.DBSubnetGroupName,
			tags: tags.TagList,
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

func (i *RDSDBSubnetGroup) Properties() types.Properties {
        properties := types.NewProperties()
        properties.Set("Name", i.name)

        for _, tag := range i.tags {
                properties.SetTag(tag.Key, tag.Value)
        }

        return properties
}
