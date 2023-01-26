package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ElasticacheReplicationGroup struct {
	svc     *elasticache.ElastiCache
	groupID *string
	tags    []*elasticache.Tag
}

func init() {
	register("ElasticacheReplicationGroup", ListElasticacheReplicationGroups)
}

func ListElasticacheReplicationGroups(sess *session.Session) ([]Resource, error) {
	svc := elasticache.New(sess)
	var resources []Resource

	params := &elasticache.DescribeReplicationGroupsInput{MaxRecords: aws.Int64(100)}

	for {
		resp, err := svc.DescribeReplicationGroups(params)
		if err != nil {
			return nil, err
		}

		for _, replicationGroup := range resp.ReplicationGroups {
			tags, err := svc.ListTagsForResource(&elasticache.ListTagsForResourceInput{
				ResourceName: replicationGroup.ARN,
			})

			if err != nil {
				continue
			}

			resources = append(resources, &ElasticacheReplicationGroup{
				svc:     svc,
				groupID: replicationGroup.ReplicationGroupId,
				tags:    tags.TagList,
			})
		}

		if resp.Marker == nil {
			break
		}

		params.Marker = resp.Marker
	}

	return resources, nil
}

func (i *ElasticacheReplicationGroup) Remove() error {
	params := &elasticache.DeleteReplicationGroupInput{
		ReplicationGroupId: i.groupID,
	}

	_, err := i.svc.DeleteReplicationGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *ElasticacheReplicationGroup) Properties() types.Properties {
	properties := types.NewProperties().
		Set("GroupID", i.groupID)

	for _, tagValue := range i.tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}

	return properties
}

func (i *ElasticacheReplicationGroup) String() string {
	return *i.groupID
}
