package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type ElasticacheReplicationGroup struct {
	svc              *elasticache.ElastiCache
	replicationGroup *elasticache.ReplicationGroup
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
			resources = append(resources, &ElasticacheReplicationGroup{
				svc:              svc,
				replicationGroup: replicationGroup,
			})
		}

		if resp.Marker == nil {
			break
		}

		params.Marker = resp.Marker
	}

	return resources, nil
}

func (i *ElasticacheReplicationGroup) Properties() types.Properties {
	properties := types.NewProperties()
	if i.replicationGroup.ReplicationGroupCreateTime != nil {
		properties.Set("CreationTime", i.replicationGroup.ReplicationGroupCreateTime.Format(time.RFC3339))
	}

	return properties
}

func (i *ElasticacheReplicationGroup) Remove() error {
	params := &elasticache.DeleteReplicationGroupInput{
		ReplicationGroupId: i.replicationGroup.ReplicationGroupId,
	}

	_, err := i.svc.DeleteReplicationGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *ElasticacheReplicationGroup) String() string {
	return *i.replicationGroup.ReplicationGroupId
}
