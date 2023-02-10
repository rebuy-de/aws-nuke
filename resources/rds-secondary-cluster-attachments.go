package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RDSSecondaryClusterAttachment struct {
	svc                *rds.RDS
	GlobalClusterID    *string
	SecondaryClusterID *string
}

func init() {
	register("RDSSecondaryClusterAttachment", ListRDSSecondaryClusterAttachment)
}

func ListRDSSecondaryClusterAttachment(sess *session.Session) ([]Resource, error) {
	svc := rds.New(sess)

	params := &rds.DescribeGlobalClustersInput{}
	resp, err := svc.DescribeGlobalClusters(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, instance := range resp.GlobalClusters {
		for _, secondaryCluster := range instance.GlobalClusterMembers {

			resources = append(resources, &RDSSecondaryClusterAttachment{
				svc:                svc,
				GlobalClusterID:    instance.GlobalClusterIdentifier,
				SecondaryClusterID: secondaryCluster.DBClusterArn,
			})
		}
	}

	return resources, nil
}

func (i *RDSSecondaryClusterAttachment) Remove() error {
	params := &rds.RemoveFromGlobalClusterInput{
		DbClusterIdentifier:     i.SecondaryClusterID,
		GlobalClusterIdentifier: i.GlobalClusterID,
	}

	_, err := i.svc.RemoveFromGlobalCluster(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *RDSSecondaryClusterAttachment) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("GlobalClusterID", i.GlobalClusterID)
	properties.Set("SecondaryClusterID", i.SecondaryClusterID)

	return properties
}

func (i *RDSSecondaryClusterAttachment) String() string {
	return aws.StringValue(i.SecondaryClusterID)
}
