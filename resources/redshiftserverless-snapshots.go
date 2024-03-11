package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshiftserverless"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RedshiftServerlessSnapshot struct {
	svc      *redshiftserverless.RedshiftServerless
	snapshot *redshiftserverless.Snapshot
}

func init() {
	register("RedshiftServerlessSnapshot", ListRedshiftServerlessSnapshots)
}

func ListRedshiftServerlessSnapshots(sess *session.Session) ([]Resource, error) {
	svc := redshiftserverless.New(sess)
	resources := []Resource{}

	params := &redshiftserverless.ListSnapshotsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListSnapshots(params)
		if err != nil {
			return nil, err
		}

		for _, snapshot := range output.Snapshots {
			resources = append(resources, &RedshiftServerlessSnapshot{
				svc:      svc,
				snapshot: snapshot,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (s *RedshiftServerlessSnapshot) Properties() types.Properties {
	properties := types.NewProperties().
		Set("CreateTime", s.snapshot.SnapshotCreateTime).
		Set("Namespace", s.snapshot.NamespaceName).
		Set("SnapshotName", s.snapshot.SnapshotName)

	return properties
}

func (s *RedshiftServerlessSnapshot) Remove() error {
	_, err := s.svc.DeleteSnapshot(&redshiftserverless.DeleteSnapshotInput{
		SnapshotName: s.snapshot.SnapshotName,
	})

	return err
}

func (s *RedshiftServerlessSnapshot) String() string {
	return *s.snapshot.SnapshotName
}
