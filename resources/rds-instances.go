package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type RDSInstance struct {
	svc                *rds.RDS
	id                 string
	deletionProtection bool
	tags		   []*rds.Tag
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
		tags, err := retrieveRDSInstanceTags(svc, *instance.DBInstanceArn)

		if err != nil {
			continue
		}

		resources = append(resources, &RDSInstance{
			svc:                svc,
			id:                 *instance.DBInstanceIdentifier,
			deletionProtection: *instance.DeletionProtection,
			tags:		    tags,
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

func retrieveRDSInstanceTags(svc *rds.RDS, instanceArn string) ([]*rds.Tag, error) {
	input := &rds.ListTagsForResourceInput{
		ResourceName: aws.String(instanceArn),
	}

	result, err := svc.ListTagsForResource(input)
	if err != nil {
		return make([]*rds.Tag, 0), err
	}

	return result.TagList, nil
}

func (i *RDSInstance) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("id", i.id)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}

func (i *RDSInstance) String() string {
	return i.id
}
