package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EKSFargateProfile struct {
	svc     *eks.EKS
	cluster *string
	name    *string
}

func init() {
	register("EKSFargateProfiles", ListEKSFargateProfiles)
}

func ListEKSFargateProfiles(sess *session.Session) ([]Resource, error) {
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

	fargateInputParams := &eks.ListFargateProfilesInput{
		MaxResults: aws.Int64(100),
	}

	// fetch the associated eks fargate profiles
	for _, clusterName := range clusterNames {
		fargateInputParams.ClusterName = clusterName

		for {
			resp, err := svc.ListFargateProfiles(fargateInputParams)
			if err != nil {
				return nil, err
			}

			for _, name := range resp.FargateProfileNames {
				resources = append(resources, &EKSFargateProfile{
					svc:     svc,
					name:    name,
					cluster: clusterName,
				})
			}

			if resp.NextToken == nil {
				fargateInputParams.NextToken = nil
				break
			}

			fargateInputParams.NextToken = resp.NextToken
		}

	}

	return resources, nil
}

func (fp *EKSFargateProfile) Remove() error {
	_, err := fp.svc.DeleteFargateProfile(&eks.DeleteFargateProfileInput{
		ClusterName:        fp.cluster,
		FargateProfileName: fp.name,
	})
	return err
}

func (fp *EKSFargateProfile) Properties() types.Properties {
	return types.NewProperties().
		Set("Cluster", *fp.cluster).
		Set("Profile", *fp.name)
}

func (fp *EKSFargateProfile) String() string {
	return fmt.Sprintf("%s:%s", *fp.cluster, *fp.name)
}
