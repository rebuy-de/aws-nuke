package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshift"
)

type RedshiftCluster struct {
	svc               *redshift.Redshift
	clusterIdentifier *string
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
				svc:               svc,
				clusterIdentifier: cluster.ClusterIdentifier,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *RedshiftCluster) Remove() error {

	_, err := f.svc.DeleteCluster(&redshift.DeleteClusterInput{
		ClusterIdentifier:        f.clusterIdentifier,
		SkipFinalClusterSnapshot: aws.Bool(true),
	})

	return err
}

func (f *RedshiftCluster) String() string {
	return *f.clusterIdentifier
}
