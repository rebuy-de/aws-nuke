package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudhsmv2"
)

type CloudHSMV2Cluster struct {
	svc       *cloudhsmv2.CloudHSMV2
	clusterID *string
}

func init() {
	register("CloudHSMV2Cluster", ListCloudHSMV2Clusters)
}

func ListCloudHSMV2Clusters(sess *session.Session) ([]Resource, error) {
	svc := cloudhsmv2.New(sess)
	resources := []Resource{}

	params := &cloudhsmv2.DescribeClustersInput{
		MaxResults: aws.Int64(25),
	}

	for {
		resp, err := svc.DescribeClusters(params)
		if err != nil {
			return nil, err
		}

		for _, cluster := range resp.Clusters {
			resources = append(resources, &CloudHSMV2Cluster{
				svc:       svc,
				clusterID: cluster.ClusterId,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *CloudHSMV2Cluster) Remove() error {

	_, err := f.svc.DeleteCluster(&cloudhsmv2.DeleteClusterInput{
		ClusterId: f.clusterID,
	})

	return err
}

func (f *CloudHSMV2Cluster) String() string {
	return *f.clusterID
}
