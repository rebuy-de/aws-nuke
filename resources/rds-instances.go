package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RDSInstance struct {
	svc                *rds.RDS
	id                 string
	deletionProtection bool
}

func init() {
	register("RDSInstance", ListRDSInstances)
}

func ListRDSInstances(sess *session.Session) ([]Resource, error) {
	svc := rds.New(sess)

	params := &rds.DescribeDBInstancesInput{}
	resp, err := svc.DescribeDBInstances(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, instance := range resp.DBInstances {
		resources = append(resources, &RDSInstance{
			svc:                svc,
			id:                 *instance.DBInstanceIdentifier,
			deletionProtection: *instance.DeletionProtection,
		})
	}

	return resources, nil
}

func (i *RDSInstance) Remove() error {
	if (i.deletionProtection) {
		modifyParams := &rds.ModifyDBInstanceInput{
			DBInstanceIdentifier: &i.id,
			DeletionProtection:   aws.Bool(false),
		}
		_, err := i.svc.ModifyDBInstance(modifyParams)
		if err != nil {
			return err
		}
	}

	params := &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier: &i.id,
		SkipFinalSnapshot:    aws.Bool(true),
	}

	_, err := i.svc.DeleteDBInstance(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *RDSInstance) String() string {
	return i.id
}
