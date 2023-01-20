package resources

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EKSNodegroup struct {
	svc       *eks.EKS
	nodegroup *eks.Nodegroup
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
	describeNodegroupInputParams := &eks.DescribeNodegroupInput{}

	// fetch the associated node groups
	for _, clusterName := range clusterNames {
		nodegroupsInputParams.ClusterName = clusterName
		describeNodegroupInputParams.ClusterName = clusterName

		for {
			resp, err := svc.ListNodegroups(nodegroupsInputParams)
			if err != nil {
				return nil, err
			}

			for _, nodegroupName := range resp.Nodegroups {
				describeNodegroupInputParams.NodegroupName = nodegroupName
				nodegroupDescriptionResponse, err := svc.DescribeNodegroup(describeNodegroupInputParams)
				if err != nil {
					return nil, err
				}
				resources = append(resources, &EKSNodegroup{
					svc:       svc,
					nodegroup: nodegroupDescriptionResponse.Nodegroup,
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
		ClusterName:   ng.nodegroup.ClusterName,
		NodegroupName: ng.nodegroup.NodegroupName,
	})
	return err
}

func (ng *EKSNodegroup) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Cluster", ng.nodegroup.ClusterName)
	properties.Set("Profile", ng.nodegroup.NodegroupName)
	if ng.nodegroup.CreatedAt != nil {
		properties.Set("CreatedAt", ng.nodegroup.CreatedAt.Format(time.RFC3339))
	}
	for k, v := range ng.nodegroup.Tags {
		properties.SetTag(&k, v)
	}
	return properties
}

func (ng *EKSNodegroup) String() string {
	return fmt.Sprintf("%s:%s", *ng.nodegroup.ClusterName, *ng.nodegroup.NodegroupName)
}
