package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

type Route53ResourceRecordSet struct {
	svc          *route53.Route53
	hostedZoneId *string
	data         *route53.ResourceRecordSet
	changeId     *string
}

func init() {
	register("Route53ResourceRecord", ListRoute53ResourceRecords)
}

func ListRoute53ResourceRecords(sess *session.Session) ([]Resource, error) {
	svc := route53.New(sess)

	resources := make([]Resource, 0)

	sub, err := ListRoute53HostedZones(sess)
	if err != nil {
		return nil, err
	}

	for _, resource := range sub {
		zone := resource.(*Route53HostedZone)
		rrs, err := ListResourceRecordsForZone(svc, zone.id)
		if err != nil {
			return nil, err
		}

		resources = append(resources, rrs...)
	}

	return resources, nil
}

func ListResourceRecordsForZone(svc *route53.Route53, zoneId *string) ([]Resource, error) {
	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId: zoneId,
	}
	resp, err := svc.ListResourceRecordSets(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, rrs := range resp.ResourceRecordSets {
		resources = append(resources, &Route53ResourceRecordSet{
			svc:          svc,
			hostedZoneId: zoneId,
			data:         rrs,
		})
	}
	return resources, nil
}

func (r *Route53ResourceRecordSet) Filter() error {
	if *r.data.Type == "NS" {
		return fmt.Errorf("cannot delete NS record")
	}

	if *r.data.Type == "SOA" {
		return fmt.Errorf("cannot delete SOA record")
	}

	return nil
}

func (r *Route53ResourceRecordSet) Remove() error {
	params := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: r.hostedZoneId,
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				&route53.Change{
					Action:            aws.String("DELETE"),
					ResourceRecordSet: r.data,
				},
			},
		},
	}

	resp, err := r.svc.ChangeResourceRecordSets(params)
	if err != nil {
		return err
	}

	r.changeId = resp.ChangeInfo.Id

	return nil
}

func (r *Route53ResourceRecordSet) String() string {
	return *r.data.Name
}
