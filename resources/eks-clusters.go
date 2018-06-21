package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
)

type EKSCluster struct {
	svc  *eks.EKS
	name *string
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
			resources = append(resources, &EKSCluster{
				svc:  svc,
				name: cluster,
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

func (f *EKSCluster) String() string {
	return *f.name
}
