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
	resources := []Resource{}
	var nextToken *string

	for {
		params := &elasticache.DescribeUserGroupsInput{
			MaxRecords: aws.Int64(100),
			Marker:     nextToken,
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

		// Check if there are more results
		if resp.Marker == nil {
			break // No more results, exit the loop
		}

		// Set the nextToken for the next iteration
		nextToken = resp.Marker
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
