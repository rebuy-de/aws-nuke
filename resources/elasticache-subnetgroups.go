package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
)

type ElasticacheSubnetGroup struct {
	svc  *elasticache.ElastiCache
	name *string
}

func init() {
	register("ElasticacheSubnetGroup", ListElasticacheSubnetGroups)
}

func ListElasticacheSubnetGroups(sess *session.Session) ([]Resource, error) {
	svc := elasticache.New(sess)

	params := &elasticache.DescribeCacheSubnetGroupsInput{MaxRecords: aws.Int64(100)}
	resp, err := svc.DescribeCacheSubnetGroups(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, subnetGroup := range resp.CacheSubnetGroups {
		resources = append(resources, &ElasticacheSubnetGroup{
			svc:  svc,
			name: subnetGroup.CacheSubnetGroupName,
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
	return *i.name
}
