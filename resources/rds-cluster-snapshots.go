package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RDSClusterSnapshot struct {
	svc      *rds.RDS
	snapshot *rds.DBClusterSnapshot
	tags     []*rds.Tag
}

func init() {
	register("RDSClusterSnapshot", ListRDSClusterSnapshots)
}

func ListRDSClusterSnapshots(sess *session.Session) ([]Resource, error) {
	svc := rds.New(sess)

	params := &rds.DescribeDBClusterSnapshotsInput{MaxRecords: aws.Int64(100)}

	resp, err := svc.DescribeDBClusterSnapshots(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, snapshot := range resp.DBClusterSnapshots {
		tags, err := svc.ListTagsForResource(&rds.ListTagsForResourceInput{
			ResourceName: snapshot.DBClusterSnapshotArn,
		})
		if err != nil {
			return nil, err
		}

		resources = append(resources, &RDSClusterSnapshot{
			svc:      svc,
			snapshot: snapshot,
			tags:     tags.TagList,
		})

	}

	return resources, nil
}

func (i *RDSClusterSnapshot) Filter() error {
	if *i.snapshot.SnapshotType == "automated" {
		return fmt.Errorf("cannot delete automated snapshots")
	}
	return nil
}

func (i *RDSClusterSnapshot) Remove() error {
	if i.snapshot.DBClusterSnapshotIdentifier == nil {
		// Sanity check to make sure the delete request does not skip the
		// identifier.
		return nil
	}

	params := &rds.DeleteDBClusterSnapshotInput{
		DBClusterSnapshotIdentifier: i.snapshot.DBClusterSnapshotIdentifier,
	}

	_, err := i.svc.DeleteDBClusterSnapshot(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *RDSClusterSnapshot) String() string {
	return *i.snapshot.DBClusterSnapshotIdentifier
}

func (i *RDSClusterSnapshot) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ARN", i.snapshot.DBClusterSnapshotArn)
	properties.Set("Identifier", i.snapshot.DBClusterSnapshotIdentifier)
	properties.Set("SnapshotType", i.snapshot.SnapshotType)
	properties.Set("Status", i.snapshot.Status)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
