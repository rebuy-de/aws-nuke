package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EKSCluster struct {
	svc     *eks.EKS
	name    *string
	cluster *eks.Cluster
}

func init() {
	register("EKSCluster", ListEKSClusters)
}

func ListEKSClusters(sess *session.Session) ([]Resource, error) {
	svc := eks.New(sess)
	resources := []Resource{}

	params := &eks.ListClustersInput{
		MaxResults: aws.Int64(100),
	}

	for {
		resp, err := svc.ListClusters(params)
		if err != nil {
			return nil, err
		}

		for _, cluster := range resp.Clusters {
			dcResp, err := svc.DescribeCluster(&eks.DescribeClusterInput{Name: cluster})
			if err != nil {
				return nil, err
			}
			resources = append(resources, &EKSCluster{
				svc:     svc,
				name:    cluster,
				cluster: dcResp.Cluster,
			})
		}
		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}
	return resources, nil
}

func (f *EKSCluster) Remove() error {

	_, err := f.svc.DeleteCluster(&eks.DeleteClusterInput{
		Name: f.name,
	})

	return err
}

func (f *EKSCluster) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("CreatedAt", f.cluster.CreatedAt.Format(time.RFC3339))
	for key, value := range f.cluster.Tags {
		properties.SetTag(&key, value)
	}
	return properties
}

func (f *EKSCluster) String() string {
	return *f.name
}
