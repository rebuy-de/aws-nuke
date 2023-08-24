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
	resources := []Resource{}
	var nextToken *string

	for {
		params := &opensearchservice.ListVpcEndpointsInput{
			NextToken: nextToken,
		}
		listResp, err := svc.ListVpcEndpoints(params)
		if err != nil {
			return nil, err
		}

		for _, vpcEndpoint := range listResp.VpcEndpointSummaryList {
			resources = append(resources, &OSVPCEndpoint{
				svc:           svc,
				vpcEndpointId: vpcEndpoint.VpcEndpointId,
			})
		}

		// Check if there are more results
		if listResp.NextToken == nil {
			break // No more results, exit the loop
		}

		// Set the nextToken for the next iteration
		nextToken = listResp.NextToken
	}

	return resources, nil
}

func (o *OSVPCEndpoint) Remove() error {
	_, err := o.svc.DeleteVpcEndpoint(&opensearchservice.DeleteVpcEndpointInput{
		VpcEndpointId: o.vpcEndpointId,
	})

	return err
}

func (o *OSVPCEndpoint) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("VpcEndpointId", o.vpcEndpointId)
	return properties
}

func (o *OSVPCEndpoint) String() string {
	return *o.vpcEndpointId
}
