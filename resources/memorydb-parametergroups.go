package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/memorydb"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type MemoryDBParameterGroup struct {
	svc    *memorydb.MemoryDB
	name   *string
	family *string
	tags   []*memorydb.Tag
}

func init() {
	register("MemoryDBParameterGroup", ListMemoryDBParameterGroups)
}

func ListMemoryDBParameterGroups(sess *session.Session) ([]Resource, error) {
	svc := memorydb.New(sess)
	var resources []Resource

	params := &memorydb.DescribeParameterGroupsInput{MaxResults: aws.Int64(100)}

	for {
		resp, err := svc.DescribeParameterGroups(params)
		if err != nil {
			return nil, err
		}

		for _, parameterGroup := range resp.ParameterGroups {
			tags, err := svc.ListTags(&memorydb.ListTagsInput{
				ResourceArn: parameterGroup.ARN,
			})

			if err != nil {
				continue
			}

			resources = append(resources, &MemoryDBParameterGroup{
				svc:    svc,
				name:   parameterGroup.Name,
				family: parameterGroup.Family,
				tags:   tags.TagList,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (i *MemoryDBParameterGroup) Filter() error {
	if strings.HasPrefix(*i.name, "default.") {
		return fmt.Errorf("Cannot delete default parameter group")
	}
	return nil
}

func (i *MemoryDBParameterGroup) Remove() error {
	params := &memorydb.DeleteParameterGroupInput{
		ParameterGroupName: i.name,
	}

	_, err := i.svc.DeleteParameterGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *MemoryDBParameterGroup) String() string {
	return *i.name
}

func (i *MemoryDBParameterGroup) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("Name", i.name).
		Set("Family", i.family)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
