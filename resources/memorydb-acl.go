package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/memorydb"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type MemoryDBACL struct {
	svc  *memorydb.MemoryDB
	name *string
	tags []*memorydb.Tag
}

func init() {
	register("MemoryDBACL", ListMemoryDBACLs)
}

func ListMemoryDBACLs(sess *session.Session) ([]Resource, error) {
	svc := memorydb.New(sess)
	var resources []Resource

	params := &memorydb.DescribeACLsInput{MaxResults: aws.Int64(50)}
	for {
		resp, err := svc.DescribeACLs(params)
		if err != nil {
			return nil, err
		}

		for _, acl := range resp.ACLs {
			tags, err := svc.ListTags(&memorydb.ListTagsInput{
				ResourceArn: acl.ARN,
			})

			if err != nil {
				continue
			}

			resources = append(resources, &MemoryDBACL{
				svc:  svc,
				name: acl.Name,
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

func (i *MemoryDBACL) Filter() error {
	if *i.name == "open-access" {
		return fmt.Errorf("open-access ACL can't be deleted")
	} else {
		return nil
	}
}

func (i *MemoryDBACL) Remove() error {
	params := &memorydb.DeleteACLInput{
		ACLName: i.name,
	}

	_, err := i.svc.DeleteACL(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *MemoryDBACL) String() string {
	return *i.name
}

func (i *MemoryDBACL) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", i.name)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
