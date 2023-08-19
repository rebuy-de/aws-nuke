package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/memorydb"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type MemoryDBCluster struct {
	svc  *memorydb.MemoryDB
	name *string
	tags []*memorydb.Tag
}

func init() {
	register("MemoryDBCluster", ListMemoryDbClusters)
}

func ListMemoryDbClusters(sess *session.Session) ([]Resource, error) {
	svc := memorydb.New(sess)
	var resources []Resource

	params := &memorydb.DescribeClustersInput{MaxResults: aws.Int64(100)}

	for {
		resp, err := svc.DescribeClusters(params)
		if err != nil {
			return nil, err
		}

		for _, cluster := range resp.Clusters {
			tags, err := svc.ListTags(&memorydb.ListTagsInput{
				ResourceArn: cluster.ARN,
			})

			if err != nil {
				continue
			}

			resources = append(resources, &MemoryDBCluster{
				svc:  svc,
				name: cluster.Name,
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

func (c *MemoryDBCluster) Remove() error {
	params := &memorydb.DeleteClusterInput{
		ClusterName: c.name,
	}

	_, err := c.svc.DeleteCluster(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *MemoryDBCluster) String() string {
	return *i.name
}

func (i *MemoryDBCluster) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", i.name)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
