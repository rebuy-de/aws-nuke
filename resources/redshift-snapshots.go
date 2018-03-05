package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshift"
)

type RedshiftSnapshot struct {
	svc                *redshift.Redshift
	snapshotIdentifier *string
}

func init() {
	register("RedshiftSnapshot", ListRedshiftSnapshots)
}

func ListRedshiftSnapshots(sess *session.Session) ([]Resource, error) {
	svc := redshift.New(sess)
	resources := []Resource{}

	params := &redshift.DescribeClusterSnapshotsInput{
		MaxRecords: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeClusterSnapshots(params)
		if err != nil {
			return nil, err
		}

		for _, snapshot := range output.Snapshots {
			resources = append(resources, &RedshiftSnapshot{
				svc:                svc,
				snapshotIdentifier: snapshot.SnapshotIdentifier,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *RedshiftSnapshot) Remove() error {

	_, err := f.svc.DeleteClusterSnapshot(&redshift.DeleteClusterSnapshotInput{
		SnapshotIdentifier: f.snapshotIdentifier,
	})

	return err
}

func (f *RedshiftSnapshot) String() string {
	return *f.snapshotIdentifier
}
