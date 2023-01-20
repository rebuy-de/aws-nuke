package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RDSDBCluster struct {
	svc                *rds.RDS
	id                 string
	deletionProtection bool
	tags               []*rds.Tag
}

func init() {
	register("RDSDBCluster", ListRDSClusters)
}

func ListRDSClusters(sess *session.Session) ([]Resource, error) {
	svc := rds.New(sess)

	params := &rds.DescribeDBClustersInput{}
	resp, err := svc.DescribeDBClusters(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, instance := range resp.DBClusters {
		tags, err := svc.ListTagsForResource(&rds.ListTagsForResourceInput{
                        ResourceName: instance.DBClusterArn,
                })

                if err != nil {
                        continue
                }

		resources = append(resources, &RDSDBCluster{
			svc:                svc,
			id:                 *instance.DBClusterIdentifier,
			deletionProtection: *instance.DeletionProtection,
			tags:               tags.TagList,
		})
	}

	return resources, nil
}

func (i *RDSDBCluster) Remove() error {
	if (i.deletionProtection) {
		modifyParams := &rds.ModifyDBClusterInput{
			DBClusterIdentifier: &i.id,
			DeletionProtection:  aws.Bool(false),
		}
		_, err := i.svc.ModifyDBCluster(modifyParams)
		if err != nil {
			return err
		}
	}

	params := &rds.DeleteDBClusterInput{
		DBClusterIdentifier: &i.id,
		SkipFinalSnapshot:   aws.Bool(true),
	}

	_, err := i.svc.DeleteDBCluster(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *RDSDBCluster) String() string {
	return i.id
}

func (i *RDSDBCluster) Properties() types.Properties {
        properties := types.NewProperties()
        properties.Set("Identifier", i.id)
	properties.Set("Deletion Protection", i.deletionProtection)

        for _, tag := range i.tags {
                properties.SetTag(tag.Key, tag.Value)
        }

        return properties
}
