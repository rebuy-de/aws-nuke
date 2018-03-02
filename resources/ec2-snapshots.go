package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2Snapshot struct {
	svc *ec2.EC2
	id  string
}

func init() {
	register("EC2Snapshot", ListEC2Snapshots)
}

func ListEC2Snapshots(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	params := &ec2.DescribeSnapshotsInput{
		OwnerIds: []*string{
			aws.String("self"),
		},
	}
	resp, err := svc.DescribeSnapshots(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.Snapshots {
		resources = append(resources, &EC2Snapshot{
			svc: svc,
			id:  *out.SnapshotId,
		})
	}

	return resources, nil
}

func (e *EC2Snapshot) Remove() error {
	_, err := e.svc.DeleteSnapshot(&ec2.DeleteSnapshotInput{
		SnapshotId: &e.id,
	})
	return err
}

func (e *EC2Snapshot) String() string {
	return e.id
}
