package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opensearchservice"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type OSVPCEndpoint struct {
	svc           *opensearchservice.OpenSearchService
	vpcEndpointId *string
}

func init() {
	register("OSVPCEndpoint", ListOSVPCEndpoints)
}

func ListOSVPCEndpoints(sess *session.Session) ([]Resource, error) {
	svc := opensearchservice.New(sess)

	vpcEndpointIds, err := getOpenSearchVpcEndpointIds(svc)
	if err != nil {
		return nil, err
	}

	listResp, err := svc.DescribeVpcEndpoints(&opensearchservice.DescribeVpcEndpointsInput{
		VpcEndpointIds: vpcEndpointIds,
	})
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	for _, vpcEndpoint := range listResp.VpcEndpoints {
		resources = append(resources, &OSVPCEndpoint{
			svc:           svc,
			vpcEndpointId: vpcEndpoint.VpcEndpointId,
		})
	}

	return resources, nil
}

func getOpenSearchVpcEndpointIds(svc *opensearchservice.OpenSearchService) ([]*string, error) {
	vpcEndpointIds := make([]*string, 0)

	listResp, err := svc.ListVpcEndpoints(&opensearchservice.ListVpcEndpointsInput{})
	if err != nil {
		return nil, err
	}

	for _, vpcEndpoint := range listResp.VpcEndpointSummaryList {
		vpcEndpointIds = append(vpcEndpointIds, vpcEndpoint.VpcEndpointId)
	}

	return vpcEndpointIds, nil
}

func (o *OSVPCEndpoint) Remove() error {
	_, err := o.svc.DeleteVpcEndpoint(&opensearchservice.DeleteVpcEndpointInput{
		VpcEndpointId: o.vpcEndpointId,
	})

	return err
}

func (o *OSVPCEndpoint) Properties() types.Properties {
	properties := types.NewProperties()
	return properties
}

func (o *OSVPCEndpoint) String() string {
	return *o.vpcEndpointId
}
