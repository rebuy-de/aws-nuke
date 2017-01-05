package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
)

type ElasticacheCacheCluster struct {
	svc       *elasticache.ElastiCache
	clusterID *string
	status    *string
}

func (n *ElasticacheNuke) ListCacheClusters() ([]Resource, error) {
	params := &elasticache.DescribeCacheClustersInput{MaxRecords: aws.Int64(100)}
	resp, err := n.Service.DescribeCacheClusters(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, cacheCluster := range resp.CacheClusters {
		resources = append(resources, &ElasticacheCacheCluster{
			svc:       n.Service,
			clusterID: cacheCluster.CacheClusterId,
			status:    cacheCluster.CacheClusterStatus,
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

func (i *ElasticacheCacheCluster) Wait() error {

	params := &elasticache.DescribeCacheClustersInput{
		CacheClusterId: i.clusterID,
	}
	return i.svc.WaitUntilCacheClusterDeleted(params)
}
