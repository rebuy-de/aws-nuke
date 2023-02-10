package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RDSGlobalCluster struct {
	svc                *rds.RDS
	id                 *string
	deletionProtection bool
}

func init() {
	register("RDSGlobalCluster", ListRDSGlobalClusters)
}

func ListRDSGlobalClusters(sess *session.Session) ([]Resource, error) {
	svc := rds.New(sess)

	params := &rds.DescribeGlobalClustersInput{}
	resp, err := svc.DescribeGlobalClusters(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, instance := range resp.GlobalClusters {
		resources = append(resources, &RDSGlobalCluster{
			svc: svc,
			id:  instance.GlobalClusterIdentifier,
		})
	}

	return resources, nil
}

func (i *RDSGlobalCluster) Remove() error {
	if i.deletionProtection {
		modifyParams := &rds.ModifyDBClusterInput{
			DBClusterIdentifier: i.id,
			DeletionProtection:  aws.Bool(false),
		}
		_, err := i.svc.ModifyDBCluster(modifyParams)
		if err != nil {
			return err
		}
	}

	params := &rds.DeleteGlobalClusterInput{
		GlobalClusterIdentifier: i.id,
	}

	_, err := i.svc.DeleteGlobalCluster(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *RDSGlobalCluster) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Identifier", i.id)
	properties.Set("DeletionProtection", i.deletionProtection)

	return properties
}

func (i *RDSGlobalCluster) String() string {
	return aws.StringValue(i.id)
}
