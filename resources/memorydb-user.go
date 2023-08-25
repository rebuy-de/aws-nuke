package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/memorydb"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type MemoryDBUser struct {
	svc  *memorydb.MemoryDB
	name *string
	tags []*memorydb.Tag
}

func init() {
	register("MemoryDBUser", ListMemoryDBUsers)
}

func ListMemoryDBUsers(sess *session.Session) ([]Resource, error) {
	svc := memorydb.New(sess)
	var resources []Resource

	params := &memorydb.DescribeUsersInput{MaxResults: aws.Int64(50)}
	for {
		resp, err := svc.DescribeUsers(params)
		if err != nil {
			return nil, err
		}

		for _, user := range resp.Users {
			tags, err := svc.ListTags(&memorydb.ListTagsInput{
				ResourceArn: user.ARN,
			})

			if err != nil {
				continue
			}

			resources = append(resources, &MemoryDBUser{
				svc:  svc,
				name: user.Name,
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

func (i *MemoryDBUser) Filter() error {
	if strings.EqualFold(*i.name, "default") {
		return fmt.Errorf("Cannot delete default user")
	}
	return nil
}

func (i *MemoryDBUser) Remove() error {
	params := &memorydb.DeleteUserInput{
		UserName: i.name,
	}

	_, err := i.svc.DeleteUser(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *MemoryDBUser) String() string {
	return *i.name
}

func (i *MemoryDBUser) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("Name", i.name)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
