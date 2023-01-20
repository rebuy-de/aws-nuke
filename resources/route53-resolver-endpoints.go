package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53resolver"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

// Route53ResolverEndpoint is the resource type for nuking
type Route53ResolverEndpoint struct {
	svc  *route53resolver.Route53Resolver
	id   *string
	name *string
}

func init() {
	register("Route53ResolverEndpoint", ListRoute53ResolverEndpoints)
}

// ListRoute53ResolverEndpoints produces the resources to be nuked
func ListRoute53ResolverEndpoints(sess *session.Session) ([]Resource, error) {
	svc := route53resolver.New(sess)

	params := &route53resolver.ListResolverEndpointsInput{}

	var resources []Resource

	for {
		resp, err := svc.ListResolverEndpoints(params)

		if err != nil {
			return nil, err
		}

		for _, endpoint := range resp.ResolverEndpoints {
			resolverEndpoint := &Route53ResolverEndpoint{
				svc:  svc,
				id:   endpoint.Id,
				name: endpoint.Name,
			}

			resources = append(resources, resolverEndpoint)
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
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
