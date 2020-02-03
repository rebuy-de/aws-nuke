package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/Optum/aws-nuke/pkg/types"
)

type EC2VPC struct {
	svc       *ec2.EC2
	vpc       *ec2.Vpc
}

func init() {
	register("EC2VPC", ListEC2VPCs)
}

func ListEC2VPCs(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeVpcs(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, vpc := range resp.Vpcs {
		resources = append(resources, &EC2VPC{
			svc:       svc,
			vpc:       vpc,
		})
	}

	return resources, nil
}

func (e *EC2VPC) Remove() error {
	params := &ec2.DeleteVpcInput{
		VpcId: e.vpc.VpcId,
	}

	_, err := e.svc.DeleteVpc(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2VPC) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range e.vpc.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	properties.Set("ID", e.vpc.VpcId)
	properties.Set("IsDefault", e.vpc.IsDefault)
	return properties
}

func (e *EC2VPC) String() string {
	return *e.vpc.VpcId
}
