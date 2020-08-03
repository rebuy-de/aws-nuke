package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

// Route53ResolverEndpoint is the resource type for nuking
type Route53ResolverEndpoint struct {
	svc  *route53resolver.Route53Resolver
	id   *string
	name *string
	ips  []*route53resolver.IpAddressUpdate
}

func init() {
	register("Route53ResolverEndpoints", ListRoute53ResolverEndpoints)
}

// ListRoute53ResolverEndpoints produces the resources to be nuked
func ListRoute53ResolverEndpoints(sess *session.Session) ([]Resource, error) {
	svc := route53resolver.New(sess)

	params := &route53resolver.ListResolverEndpointsInput{}

	resources := make([]Resource, 0)
	output, err := svc.ListResolverEndpoints(params)

	if err != nil {
		return resources, err
	}

	for _, endpoint := range output.ResolverEndpoints {
		resolverEndpoint := &Route53ResolverEndpoint{
			svc:  svc,
			id:   endpoint.Id,
			name: endpoint.Name,
		}

		ipsOutput, err := svc.ListResolverEndpointIpAddresses(
			&route53resolver.ListResolverEndpointIpAddressesInput{
				ResolverEndpointId: endpoint.Id,
			})

		if err != nil {
			return resources, err
		}

		for _, ip := range ipsOutput.IpAddresses {
			resolverEndpoint.ips = append(resolverEndpoint.ips, &route53resolver.IpAddressUpdate{
				Ip:       ip.Ip,
				IpId:     ip.IpId,
				SubnetId: ip.SubnetId,
			})
		}

		resources = append(resources, resolverEndpoint)
	}

	return resources, err
}

// Remove implements Resource
func (r *Route53ResolverEndpoint) Remove() error {
	_, err := r.svc.DeleteResolverEndpoint(
		&route53resolver.DeleteResolverEndpointInput{
			ResolverEndpointId: r.id,
		})

	if err != nil {
		return err
	}

	return nil
}

// Properties provides debugging output
func (r *Route53ResolverEndpoint) Properties() types.Properties {
	return types.NewProperties().
		Set("EndpointID", r.id).
		Set("Name", r.name)
}

// String implements Stringer
func (r *Route53ResolverEndpoint) String() string {
	return fmt.Sprintf("%s (%s)", *r.id, *r.name)
}
