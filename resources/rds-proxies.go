package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RDSProxy struct {
	svc  *rds.RDS
	id   string
	tags []*rds.Tag
}

func init() {
	register("RDSProxy", ListRDSProxies)
}

func ListRDSProxies(sess *session.Session) ([]Resource, error) {
	svc := rds.New(sess)

	params := &rds.DescribeDBProxiesInput{}
	resp, err := svc.DescribeDBProxies(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, instance := range resp.DBProxies {
		tags, err := svc.ListTagsForResource(&rds.ListTagsForResourceInput{
			ResourceName: instance.DBProxyArn,
		})

		if err != nil {
			continue
		}

		resources = append(resources, &RDSProxy{
			svc:  svc,
			id:   *instance.DBProxyName,
			tags: tags.TagList,
		})
	}

	return resources, nil
}

func (i *RDSProxy) Remove() error {
	params := &rds.DeleteDBProxyInput{
		DBProxyName: &i.id,
	}

	_, err := i.svc.DeleteDBProxy(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *RDSProxy) String() string {
	return i.id
}

func (i *RDSProxy) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ProxyName", i.id)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
