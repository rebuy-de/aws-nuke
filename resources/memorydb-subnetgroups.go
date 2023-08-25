package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/memorydb"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type MemoryDBSubnetGroup struct {
	svc  *memorydb.MemoryDB
	name *string
	tags []*memorydb.Tag
}

func init() {
	register("MemoryDBSubnetGroup", ListMemoryDBSubnetGroups)
}

func ListMemoryDBSubnetGroups(sess *session.Session) ([]Resource, error) {
	svc := memorydb.New(sess)
	var resources []Resource

	params := &memorydb.DescribeSubnetGroupsInput{MaxResults: aws.Int64(100)}

	for {
		resp, err := svc.DescribeSubnetGroups(params)
		if err != nil {
			return nil, err
		}
		for _, subnetGroup := range resp.SubnetGroups {
			tags, err := svc.ListTags(&memorydb.ListTagsInput{
				ResourceArn: subnetGroup.ARN,
			})

			if err != nil {
				continue
			}

			resources = append(resources, &MemoryDBSubnetGroup{
				svc:  svc,
				name: subnetGroup.Name,
				tags: tags.TagList,
			})

		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (i *MemoryDBSubnetGroup) Remove() error {
	params := &memorydb.DeleteSubnetGroupInput{
		SubnetGroupName: i.name,
	}

	_, err := i.svc.DeleteSubnetGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *MemoryDBSubnetGroup) String() string {
	return *i.name
}

func (i *MemoryDBSubnetGroup) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("Name", i.name)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
