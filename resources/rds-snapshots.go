package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RDSSnapshot struct {
	svc      *rds.RDS
	snapshot *rds.DBSnapshot
	tags     []*rds.Tag
}

func init() {
	register("RDSSnapshot", ListRDSSnapshots)
}

func ListRDSSnapshots(sess *session.Session) ([]Resource, error) {
	svc := rds.New(sess)

	params := &rds.DescribeDBSnapshotsInput{MaxRecords: aws.Int64(100)}
	resp, err := svc.DescribeDBSnapshots(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, snapshot := range resp.DBSnapshots {
		tags, err := svc.ListTagsForResource(&rds.ListTagsForResourceInput{
			ResourceName: snapshot.DBSnapshotArn,
		})
		if err != nil {
			return nil, err
		}

		resources = append(resources, &RDSSnapshot{
			svc:      svc,
			snapshot: snapshot,
			tags:     tags.TagList,
		})

	}

	return resources, nil
}

func (i *RDSSnapshot) Filter() error {
	if *i.snapshot.SnapshotType == "automated" {
		return fmt.Errorf("cannot delete automated snapshots")
	}
	return nil
}

func (i *RDSSnapshot) Remove() error {
	if i.snapshot.DBSnapshotIdentifier == nil {
		// Sanity check to make sure the delete request does not skip the
		// identifier.
		return nil
	}

	params := &rds.DeleteDBSnapshotInput{
		DBSnapshotIdentifier: i.snapshot.DBSnapshotIdentifier,
	}

	_, err := i.svc.DeleteDBSnapshot(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *RDSSnapshot) String() string {
	return *i.snapshot.DBSnapshotIdentifier
}

func (i *RDSSnapshot) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ARN", i.snapshot.DBSnapshotArn)
	properties.Set("Identifier", i.snapshot.DBSnapshotIdentifier)
	properties.Set("SnapshotType", i.snapshot.SnapshotType)
	properties.Set("Status", i.snapshot.Status)
	properties.Set("AvailabilityZone", i.snapshot.AvailabilityZone)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
