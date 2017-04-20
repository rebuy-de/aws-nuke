package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
)

type ElasticacheSubnetGroup struct {
	svc    *elasticache.ElastiCache
	name   *string
	region *string
}

func (n *ElasticacheNuke) ListSubnetGroups() ([]Resource, error) {
	params := &elasticache.DescribeCacheSubnetGroupsInput{MaxRecords: aws.Int64(100)}
	resp, err := n.Service.DescribeCacheSubnetGroups(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, subnetGroup := range resp.CacheSubnetGroups {
		resources = append(resources, &ElasticacheSubnetGroup{
			svc:    n.Service,
			name:   subnetGroup.CacheSubnetGroupName,
			region: n.Service.Config.Region,
		})

	}

	return resources, nil
}

func (i *ElasticacheSubnetGroup) Remove() error {
	params := &elasticache.DeleteCacheSubnetGroupInput{
		CacheSubnetGroupName: i.name,
	}

	_, err := i.svc.DeleteCacheSubnetGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *ElasticacheSubnetGroup) String() string {
	return fmt.Sprintf("%s in %s", *i.name, *i.region)
}
