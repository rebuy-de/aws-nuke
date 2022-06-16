package resources

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/config"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2Host struct {
	svc  *ec2.EC2
	host *ec2.Host

	featureFlags config.FeatureFlags
}

func init() {
	register("EC2Host", ListEC2Hosts)
}

func ListEC2Hosts(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	params := &ec2.DescribeHostsInput{}
	resources := make([]Resource, 0)
	for {
		resp, err := svc.DescribeHosts(params)
		if err != nil {
			return nil, err
		}

		for _, host := range resp.Hosts {
			resources = append(resources, &EC2Host{
				svc:  svc,
				host: host,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params = &ec2.DescribeHostsInput{
			NextToken: resp.NextToken,
		}
	}

	return resources, nil
}

func (i *EC2Host) Filter() error {
	if *i.host.State == "released" {
		return fmt.Errorf("already released")
	}
	return nil
}

func (i *EC2Host) Remove() error {
	params := &ec2.ReleaseHostsInput{
		HostIds: []*string{i.host.HostId},
	}

	_, err := i.svc.ReleaseHosts(params)
	if err != nil {
		return err
	}
	return nil
}

func (i *EC2Host) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Identifier", i.host.HostId)
	properties.Set("HostInstanceFamily", i.host.HostProperties.InstanceFamily)
	properties.Set("HostCores", i.host.HostProperties.Cores)
	properties.Set("HostState", i.host.State)
	properties.Set("AllocationTime", i.host.AllocationTime.Format(time.RFC3339))

	for _, tagValue := range i.host.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}

	return properties
}

func (i *EC2Host) String() string {
	return *i.host.HostId
}
