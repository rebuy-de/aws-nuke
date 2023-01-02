package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RedshiftSnapshot struct {
	svc      *redshift.Redshift
	snapshot *redshift.Snapshot
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
				svc:      svc,
				snapshot: snapshot,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *RedshiftSnapshot) Properties() types.Properties {
	properties := types.NewProperties().
		Set("CreatedTime", f.snapshot.SnapshotCreateTime)

	for _, tag := range f.snapshot.Tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}

func (f *RedshiftSnapshot) Remove() error {

	_, err := f.svc.DeleteClusterSnapshot(&redshift.DeleteClusterSnapshotInput{
		SnapshotIdentifier: f.snapshot.SnapshotIdentifier,
	})

	return err
}

func (f *RedshiftSnapshot) String() string {
	return *f.snapshot.SnapshotIdentifier
}
