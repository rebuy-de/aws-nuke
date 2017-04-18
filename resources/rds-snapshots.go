package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RDSSnapshot struct {
	svc        *rds.RDS
	identifier *string
	status     *string
	region     *string
}

func (n *RDSNuke) ListSnapshots() ([]Resource, error) {
	params := &rds.DescribeDBSnapshotsInput{MaxRecords: aws.Int64(100)}
	resp, err := n.Service.DescribeDBSnapshots(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, snapshot := range resp.DBSnapshots {
		resources = append(resources, &RDSSnapshot{
			svc:        n.Service,
			identifier: snapshot.DBSnapshotIdentifier,
			status:     snapshot.Status,
			region:     n.Service.Config.Region,
		})

	}

	return resources, nil
}
func (i *RDSSnapshot) Remove() error {
	params := &rds.DeleteDBSnapshotInput{
		DBSnapshotIdentifier: i.identifier,
	}

	_, err := i.svc.DeleteDBSnapshot(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *RDSSnapshot) String() string {
	return fmt.Sprintf("%s in %s", *i.identifier, *i.region)
}
