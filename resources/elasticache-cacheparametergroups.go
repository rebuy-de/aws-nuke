package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ElasticacheCacheParameterGroup struct {
	svc         *elasticache.ElastiCache
	groupName   *string
	groupFamily *string
}

func init() {
	register("ElasticacheCacheParameterGroup", ListElasticacheCacheParameterGroups)
}

func ListElasticacheCacheParameterGroups(sess *session.Session) ([]Resource, error) {
	svc := elasticache.New(sess)
	var resources []Resource

	params := &elasticache.DescribeCacheParameterGroupsInput{MaxRecords: aws.Int64(100)}

	for {
		resp, err := svc.DescribeCacheParameterGroups(params)
		if err != nil {
			return nil, err
		}

		for _, cacheParameterGroup := range resp.CacheParameterGroups {
			resources = append(resources, &ElasticacheCacheParameterGroup{
				svc:         svc,
				groupName:   cacheParameterGroup.CacheParameterGroupName,
				groupFamily: cacheParameterGroup.CacheParameterGroupFamily,
			})
		}

		if resp.Marker == nil {
			break
		}

		params.Marker = resp.Marker
	}

	return resources, nil
}

func (i *ElasticacheCacheParameterGroup) Filter() error {
	if strings.HasPrefix(*i.groupName, "default.") {
		return fmt.Errorf("Cannot delete default cache parameter group")
	}
	return nil
}

func (i *ElasticacheCacheParameterGroup) Remove() error {
	params := &elasticache.DeleteCacheParameterGroupInput{
		CacheParameterGroupName: i.groupName,
	}

	_, err := i.svc.DeleteCacheParameterGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *ElasticacheCacheParameterGroup) String() string {
	return *i.groupName
}

func (i *ElasticacheCacheParameterGroup) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("GroupName", i.groupName).
		Set("GroupFamily", i.groupFamily)
	return properties
}
