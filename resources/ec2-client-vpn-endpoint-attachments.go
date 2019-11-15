package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2ClientVpnEndpointAttachments struct {
	svc                 *ec2.EC2
	associationId       *string
	clientVpnEndpointId *string
	vpcId               *string
}

func init() {
	register("EC2ClientVpnEndpointAttachment", ListEC2ClientVpnEndpointAttachments)
}

func ListEC2ClientVpnEndpointAttachments(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	endpoints := make([]*string, 0)

	params := &ec2.DescribeClientVpnEndpointsInput{}
	err := svc.DescribeClientVpnEndpointsPages(params,
		func(page *ec2.DescribeClientVpnEndpointsOutput, lastPage bool) bool {
			for _, out := range page.ClientVpnEndpoints {
				endpoints = append(endpoints, out.ClientVpnEndpointId)
			}
			return true
		})
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, clientVpnEndpointId := range endpoints {
		params := &ec2.DescribeClientVpnTargetNetworksInput{
			ClientVpnEndpointId: clientVpnEndpointId,
		}
		err := svc.DescribeClientVpnTargetNetworksPages(params,
			func(page *ec2.DescribeClientVpnTargetNetworksOutput, lastPage bool) bool {
				for _, out := range page.ClientVpnTargetNetworks {
					resources = append(resources, &EC2ClientVpnEndpointAttachments{
						svc:                 svc,
						associationId:       out.AssociationId,
						clientVpnEndpointId: out.ClientVpnEndpointId,
						vpcId:               out.VpcId,
					})
				}
				return true
			})
		if err != nil {
			return nil, err
		}
	}

	return resources, nil
}

func (e *EC2ClientVpnEndpointAttachments) Remove() error {
	params := &ec2.DisassociateClientVpnTargetNetworkInput{
		AssociationId:       e.associationId,
		ClientVpnEndpointId: e.clientVpnEndpointId,
	}

	_, err := e.svc.DisassociateClientVpnTargetNetwork(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *EC2ClientVpnEndpointAttachments) String() string {
	return fmt.Sprintf("%s -> %s", *e.clientVpnEndpointId, *e.vpcId)
}
