package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ECSCluster struct {
	svc  *ecs.ECS
	ARN  *string
	tags []*ecs.Tag
}

func init() {
	register("ECSCluster", ListECSClusters)
}

func ListECSClusters(sess *session.Session) ([]Resource, error) {
	svc := ecs.New(sess)
	resources := []Resource{}

	params := &ecs.ListClustersInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListClusters(params)
		if err != nil {
			return nil, err
		}

		for _, clusterArn := range output.ClusterArns {
			tagsResp, err := svc.ListTagsForResource(&ecs.ListTagsForResourceInput{
				ResourceArn: clusterArn,
			})
			if err != nil {
				return nil, err
			}

			resources = append(resources, &ECSCluster{
				svc:  svc,
				ARN:  clusterArn,
				tags: tagsResp.Tags,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *ECSCluster) Remove() error {

	_, err := f.svc.DeleteCluster(&ecs.DeleteClusterInput{
		Cluster: f.ARN,
	})

	return err
}

func (f *ECSCluster) Properties() types.Properties {
	properties := types.NewProperties()

	for _, tag := range f.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}

func (f *ECSCluster) String() string {
	return *f.ARN
}
