package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/neptune"
)

type NeptuneCluster struct {
	svc *neptune.Neptune
	ID  *string
}

func init() {
	register("NeptuneCluster", ListNeptuneClusters)
}

func ListNeptuneClusters(sess *session.Session) ([]Resource, error) {
	svc := neptune.New(sess)
	resources := []Resource{}

	params := &neptune.DescribeDBClustersInput{
		MaxRecords: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeDBClusters(params)
		if err != nil {
			return nil, err
		}

		for _, dbCluster := range output.DBClusters {
			resources = append(resources, &NeptuneCluster{
				svc: svc,
				ID:  dbCluster.DBClusterIdentifier,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *NeptuneCluster) Remove() error {

	_, err := f.svc.DeleteDBCluster(&neptune.DeleteDBClusterInput{
		DBClusterIdentifier: f.ID,
		SkipFinalSnapshot:   aws.Bool(true),
	})

	return err
}

func (f *NeptuneCluster) String() string {
	return *f.ID
}
