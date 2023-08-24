package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2InstanceConnectEndpoint struct {
	svc         *ec2.EC2
	az          *string
	createdAt   *time.Time
	dnsName     *string
	fipsDNSName *string
	id          *string
	ownerID     *string
	state       *string
	subnetID    *string
	tags        []*ec2.Tag
	vpcID       *string
}

func init() {
	register("EC2InstanceConnectEndpoint", ListEC2InstanceConnectEndpoints)
}

func ListEC2InstanceConnectEndpoints(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	params := &ec2.DescribeInstanceConnectEndpointsInput{}
	resources := make([]Resource, 0)
	for {
		resp, err := svc.DescribeInstanceConnectEndpoints(params)
		if err != nil {
			return nil, err
		}

		for _, endpoint := range resp.InstanceConnectEndpoints {
			resources = append(resources, &EC2InstanceConnectEndpoint{
				svc:         svc,
				az:          endpoint.AvailabilityZone,
				createdAt:   endpoint.CreatedAt,
				dnsName:     endpoint.DnsName,
				fipsDNSName: endpoint.FipsDnsName,
				id:          endpoint.InstanceConnectEndpointId,
				ownerID:     endpoint.OwnerId,
				state:       endpoint.State,
				subnetID:    endpoint.SubnetId,
				tags:        endpoint.Tags,
				vpcID:       endpoint.VpcId,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (i *EC2InstanceConnectEndpoint) Remove() error {
	params := &ec2.DeleteInstanceConnectEndpointInput{
		InstanceConnectEndpointId: i.id,
	}

	_, err := i.svc.DeleteInstanceConnectEndpoint(params)
	if err != nil {
		return err
	}
	return nil
}

func (i *EC2InstanceConnectEndpoint) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ID", i.id)
	properties.Set("AZ", i.az)
	properties.Set("CreatedAt", i.createdAt.Format(time.RFC3339))
	properties.Set("DNSName", i.dnsName)
	properties.Set("FIPSDNSName", i.fipsDNSName)
	properties.Set("OwnerID", i.ownerID)
	properties.Set("State", i.state)
	properties.Set("SubnetID", i.subnetID)
	properties.Set("VPCID", i.vpcID)

	for _, tagValue := range i.tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}

	return properties
}

func (i *EC2InstanceConnectEndpoint) String() string {
	return *i.id
}
