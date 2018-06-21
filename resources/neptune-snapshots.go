package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/neptune"
)

type NetpuneSnapshot struct {
	svc *neptune.Neptune
	ID  *string
}

func init() {
	register("NetpuneSnapshot", ListNetpuneSnapshots)
}

func ListNetpuneSnapshots(sess *session.Session) ([]Resource, error) {
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
			resources = append(resources, &NetpuneSnapshot{
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

func (f *NetpuneSnapshot) Remove() error {

	_, err := f.svc.DeleteDBClusterSnapshot(&neptune.DeleteDBClusterSnapshotInput{
		DBClusterSnapshotIdentifier: f.ID,
	})

	return err
}

func (f *NetpuneSnapshot) String() string {
	return *f.ID
}
