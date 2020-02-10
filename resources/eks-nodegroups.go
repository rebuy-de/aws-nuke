package resources

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type EKSNodegroup struct {
	svc     *eks.EKS
	cluster *string
	name    *string
}

func init() {
	register("EKSNodegroups", ListEKSNodegroups)
}

func ListEKSNodegroups(sess *session.Session) ([]Resource, error) {
	svc := eks.New(sess)
	clusterNames := []*string{}
	resources := []Resource{}

	clusterInputParams := &eks.ListClustersInput{
		MaxResults: aws.Int64(100),
	}

	// fetch all cluster names
	for {
		resp, err := svc.ListClusters(clusterInputParams)
		if err != nil {
			return nil, err
		}

		for _, cluster := range resp.Clusters {
			clusterNames = append(clusterNames, cluster)
		}

		if resp.NextToken == nil {
			break
		}

		clusterInputParams.NextToken = resp.NextToken
	}

	nodegroupsInputParams := &eks.ListNodegroupsInput{
		MaxResults: aws.Int64(100),
	}

	// fetch the associated node groups
	for _, clusterName := range clusterNames {
		nodegroupsInputParams.ClusterName = clusterName

		for {
			resp, err := svc.ListNodegroups(nodegroupsInputParams)
			if err != nil {
				return nil, err
			}

			for _, name := range resp.Nodegroups {
				resources = append(resources, &EKSNodegroup{
					svc:     svc,
					name:    name,
					cluster: clusterName,
				})
			}

			if resp.NextToken == nil {
				nodegroupsInputParams.NextToken = nil
				break
			}

			nodegroupsInputParams.NextToken = resp.NextToken
		}

	}

	return resources, nil
}

func (ng *EKSNodegroup) Remove() error {
	_, err := ng.svc.DeleteNodegroup(&eks.DeleteNodegroupInput{
		ClusterName:   ng.cluster,
		NodegroupName: ng.name,
	})
	return err
}

func (ng *EKSNodegroup) Properties() types.Properties {
	return types.NewProperties().
		Set("Cluster", *ng.cluster).
		Set("Profile", *ng.name)
}

func (ng *EKSNodegroup) String() string {
	return fmt.Sprintf("%s:%s", *ng.cluster, *ng.name)
}
