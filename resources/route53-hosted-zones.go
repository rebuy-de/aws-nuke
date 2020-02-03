package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/Optum/aws-nuke/pkg/types"
)

func init() {
	register("Route53HostedZone", ListRoute53HostedZones)
}

func ListRoute53HostedZones(sess *session.Session) ([]Resource, error) {
	svc := route53.New(sess)

	params := &route53.ListHostedZonesInput{}
	resp, err := svc.ListHostedZones(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, hz := range resp.HostedZones {
		resources = append(resources, &Route53HostedZone{
			svc:  svc,
			id:   hz.Id,
			name: hz.Name,
		})
	}
	return resources, nil
}

type Route53HostedZone struct {
	svc  *route53.Route53
	id   *string
	name *string
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

func (hz *Route53HostedZone) Properties() types.Properties {
	return types.NewProperties().
		Set("Name", hz.name)
}

func (hz *Route53HostedZone) String() string {
	return fmt.Sprintf("%s (%s)", *hz.id, *hz.name)
}
