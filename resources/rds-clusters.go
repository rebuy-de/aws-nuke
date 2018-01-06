package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RDSDBCluster struct {
	svc *rds.RDS
	id  string
}

func init() {
	register("RDSCluster", ListRDSClusters)
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
		resources = append(resources, &RDSDBCluster{
			svc: svc,
			id:  *instance.DBClusterIdentifier,
		})
	}

	return resources, nil
}

func (i *RDSDBCluster) Remove() error {
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
