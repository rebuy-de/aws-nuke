package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/neptune"
)

type NeptuneSnapshot struct {
	svc *neptune.Neptune
	ID  *string
}

func init() {
	register("NeptuneSnapshot", ListNeptuneSnapshots)
}

func ListNeptuneSnapshots(sess *session.Session) ([]Resource, error) {
	svc := neptune.New(sess)
	resources := []Resource{}

	params := &neptune.DescribeDBClusterSnapshotsInput{
		MaxRecords: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeDBClusterSnapshots(params)
		if err != nil {
			return nil, err
		}

		for _, dbClusterSnapshot := range output.DBClusterSnapshots {
			resources = append(resources, &NeptuneSnapshot{
				svc: svc,
				ID:  dbClusterSnapshot.DBClusterSnapshotIdentifier,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *NeptuneSnapshot) Remove() error {

	_, err := f.svc.DeleteDBClusterSnapshot(&neptune.DeleteDBClusterSnapshotInput{
		DBClusterSnapshotIdentifier: f.ID,
	})

	return err
}

func (f *NeptuneSnapshot) String() string {
	return *f.ID
}
