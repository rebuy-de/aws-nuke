package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RedshiftCluster struct {
	svc     *redshift.Redshift
	cluster *redshift.Cluster
}

func init() {
	register("RedshiftCluster", ListRedshiftClusters)
}

func ListRedshiftClusters(sess *session.Session) ([]Resource, error) {
	svc := redshift.New(sess)
	resources := []Resource{}

	params := &redshift.DescribeClustersInput{
		MaxRecords: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeClusters(params)
		if err != nil {
			return nil, err
		}

		for _, cluster := range output.Clusters {
			resources = append(resources, &RedshiftCluster{
				svc:     svc,
				cluster: cluster,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *RedshiftCluster) Properties() types.Properties {
	properties := types.NewProperties().
		Set("CreatedTime", f.cluster.ClusterCreateTime)

	for _, tag := range f.cluster.Tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}

func (f *RedshiftCluster) Remove() error {

	_, err := f.svc.DeleteCluster(&redshift.DeleteClusterInput{
		ClusterIdentifier:        f.cluster.ClusterIdentifier,
		SkipFinalClusterSnapshot: aws.Bool(true),
	})

	return err
}

func (f *RedshiftCluster) String() string {
	return *f.cluster.ClusterIdentifier
}
