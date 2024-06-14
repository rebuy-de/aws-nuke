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
	profile *eks.FargateProfile
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

		clusterNames = append(clusterNames, resp.Clusters...)

		if resp.NextToken == nil {
			break
		}

		clusterInputParams.NextToken = resp.NextToken
	}

	listFargateInputParams := &eks.ListFargateProfilesInput{
		MaxResults: aws.Int64(100),
	}
	describeFargateInputParams := &eks.DescribeFargateProfileInput{}

	// fetch the associated eks fargate profiles
	for _, clusterName := range clusterNames {
		listFargateInputParams.ClusterName = clusterName
		describeFargateInputParams.ClusterName = clusterName

		for {
			resp, err := svc.ListFargateProfiles(listFargateInputParams)
			if err != nil {
				return nil, err
			}

			for _, name := range resp.FargateProfileNames {
				describeFargateInputParams.FargateProfileName = name
				describeFargateResp, err := svc.DescribeFargateProfile(describeFargateInputParams)
				if err != nil {
					return nil, err
				}
				resources = append(resources, &EKSFargateProfile{
					svc:     svc,
					profile: describeFargateResp.FargateProfile,
				})
			}

			if resp.NextToken == nil {
				listFargateInputParams.NextToken = nil
				break
			}

			listFargateInputParams.NextToken = resp.NextToken
		}

	}

	return resources, nil
}

func (fp *EKSFargateProfile) Remove() error {
	_, err := fp.svc.DeleteFargateProfile(&eks.DeleteFargateProfileInput{
		ClusterName:        fp.profile.ClusterName,
		FargateProfileName: fp.profile.FargateProfileName,
	})
	return err
}

func (fp *EKSFargateProfile) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Cluster", fp.profile.ClusterName)
	properties.Set("Profile", fp.profile.FargateProfileName)
	for k, v := range fp.profile.Tags {
		properties.SetTag(&k, v)
	}
	return properties
}

func (fp *EKSFargateProfile) String() string {
	return fmt.Sprintf("%s:%s", *fp.profile.ClusterName, *fp.profile.FargateProfileName)
}
