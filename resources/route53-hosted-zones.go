package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/route53"
)

func (n *Route53Nuke) ListHostedZones() ([]Resource, error) {
	params := &route53.ListHostedZonesInput{}
	resp, err := n.Service.ListHostedZones(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, hz := range resp.HostedZones {
		resources = append(resources, &Route53HostedZone{
			svc:    n.Service,
			id:     hz.Id,
			name:   hz.Name,
			region: n.Service.Config.Region,
		})
	}
	return resources, nil
}

type Route53HostedZone struct {
	svc    *route53.Route53
	id     *string
	name   *string
	region *string
}

func (hz *Route53HostedZone) Remove() error {
	params := &route53.DeleteHostedZoneInput{
		Id: hz.id,
	}

	_, err := hz.svc.DeleteHostedZone(params)
	if err != nil {
		return err
	}

	return nil
}

func (hz *Route53HostedZone) String() string {
	return fmt.Sprintf("%s (%s) in %s", *hz.id, *hz.name, *hz.region)
}
