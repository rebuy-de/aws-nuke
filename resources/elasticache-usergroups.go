package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
)

type ElasticacheUserGroup struct {
	svc     *elasticache.ElastiCache
	groupId *string
}

func init() {
	register("ElasticacheUserGroup", ListElasticacheUserGroups)
}

func ListElasticacheUserGroups(sess *session.Session) ([]Resource, error) {
	svc := elasticache.New(sess)

	var resources []Resource
	var marker *string // Marker for pagination

	for {
		params := &elasticache.DescribeUserGroupsInput{
			MaxRecords: aws.Int64(100),
			Marker:     marker,
		}
		resp, err := svc.DescribeUserGroups(params)
		if err != nil {
			return nil, err
		}

		for _, userGroup := range resp.UserGroups {
			resources = append(resources, &ElasticacheUserGroup{
				svc:     svc,
				groupId: userGroup.UserGroupId,
			})
		}

		// If there are more results, the response will have a Marker.
		// Set the marker for the next iteration.
		if resp.Marker != nil {
			marker = resp.Marker
		} else {
			// No more results, break the loop.
			break
		}
	}

	return resources, nil
}

func (i *ElasticacheUserGroup) Remove() error {
	params := &elasticache.DeleteUserGroupInput{
		UserGroupId: i.groupId,
	}

	_, err := i.svc.DeleteUserGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *ElasticacheUserGroup) String() string {
	return *i.groupId
}
