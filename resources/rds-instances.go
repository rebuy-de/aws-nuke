package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/rebuy-de/aws-nuke/v2/pkg/config"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RDSInstance struct {
	svc      *rds.RDS
	instance *rds.DBInstance
	tags     []*rds.Tag

	featureFlags config.FeatureFlags
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
		tags, err := svc.ListTagsForResource(&rds.ListTagsForResourceInput{
			ResourceName: instance.DBInstanceArn,
		})

		if err != nil {
			continue
		}

		resources = append(resources, &RDSInstance{
			svc:      svc,
			instance: instance,
			tags:     tags.TagList,
		})
	}

	return resources, nil
}

func (i *RDSInstance) FeatureFlags(ff config.FeatureFlags) {
	i.featureFlags = ff
}

func (i *RDSInstance) Remove() error {
	if aws.BoolValue(i.instance.DeletionProtection) && i.featureFlags.DisableDeletionProtection.RDSInstance {
		modifyParams := &rds.ModifyDBInstanceInput{
			DBInstanceIdentifier: i.instance.DBInstanceIdentifier,
			DeletionProtection:   aws.Bool(false),
		}
		_, err := i.svc.ModifyDBInstance(modifyParams)
		if err != nil {
			return err
		}
	}

	params := &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier: i.instance.DBInstanceIdentifier,
		SkipFinalSnapshot:    aws.Bool(true),
	}

	_, err := i.svc.DeleteDBInstance(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *RDSInstance) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Identifier", i.instance.DBInstanceIdentifier)
	properties.Set("DeletionProtection", i.instance.DeletionProtection)
	properties.Set("AvailabilityZone", i.instance.AvailabilityZone)
	properties.Set("InstanceClass", i.instance.DBInstanceClass)
	if i.instance.InstanceCreateTime != nil {
		properties.Set("InstanceCreateTime", i.instance.InstanceCreateTime.Format(time.RFC3339))
	}
	properties.Set("Engine", i.instance.Engine)
	properties.Set("EngineVersion", i.instance.EngineVersion)
	properties.Set("MultiAZ", i.instance.MultiAZ)
	properties.Set("PubliclyAccessible", i.instance.PubliclyAccessible)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}

func (i *RDSInstance) String() string {
	return aws.StringValue(i.instance.DBInstanceIdentifier)
}
