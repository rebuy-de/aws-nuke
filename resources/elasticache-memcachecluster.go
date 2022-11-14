package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ElasticacheCacheCluster struct {
	svc       *elasticache.ElastiCache
	clusterID *string
	status    *string
	tags      []*elasticache.Tag
}

func init() {
	register("ElasticacheCacheCluster", ListElasticacheCacheClusters)
}

func ListElasticacheCacheClusters(sess *session.Session) ([]Resource, error) {
	svc := elasticache.New(sess)

	// Lookup current account ID
	stsSvc := sts.New(sess)
	callerID, err := stsSvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}
	accountID := callerID.Account
	region := svc.Config.Region

	params := &elasticache.DescribeCacheClustersInput{MaxRecords: aws.Int64(100)}
	resp, err := svc.DescribeCacheClusters(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, cacheCluster := range resp.CacheClusters {
		// Arn creation for listing tags
		tags, err := svc.ListTagsForResource(&elasticache.ListTagsForResourceInput{
			ResourceName: aws.String(fmt.Sprintf("arn:aws:elasticache:%s:%s:cluster:%s", *region, *accountID, *cacheCluster.CacheClusterId)),
		})
		if err != nil {
			continue
		}
		resources = append(resources, &ElasticacheCacheCluster{
			svc:       svc,
			clusterID: cacheCluster.CacheClusterId,
			status:    cacheCluster.CacheClusterStatus,
			tags:      tags.TagList,
		})

	}
	return resources, nil
}

func (i *ElasticacheCacheCluster) Remove() error {
	params := &elasticache.DeleteCacheClusterInput{
		CacheClusterId: i.clusterID,
	}

	_, err := i.svc.DeleteCacheCluster(params)
	if err != nil {
		return err
	}
	return nil
}

func (i *ElasticacheCacheCluster) String() string {
	return *i.clusterID
}

func (i *ElasticacheCacheCluster) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Identifier", i.clusterID)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
