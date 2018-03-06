package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dax"
)

type DAXCluster struct {
	svc         *dax.DAX
	clusterName *string
}

func init() {
	register("DAXCluster", ListDAXClusters)
}

func ListDAXClusters(sess *session.Session) ([]Resource, error) {
	svc := dax.New(sess)
	resources := []Resource{}

	params := &dax.DescribeClustersInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeClusters(params)
		if err != nil {
			return nil, err
		}

		for _, cluster := range output.Clusters {
			resources = append(resources, &DAXCluster{
				svc:         svc,
				clusterName: cluster.ClusterName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *DAXCluster) Remove() error {

	_, err := f.svc.DeleteCluster(&dax.DeleteClusterInput{
		ClusterName: f.clusterName,
	})

	return err
}

func (f *DAXCluster) String() string {
	return *f.clusterName
}
