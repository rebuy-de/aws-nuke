package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type Route53ResourceRecordSet struct {
	svc            *route53.Route53
	hostedZoneId   *string
	hostedZoneName *string
	data           *route53.ResourceRecordSet
	changeId       *string
	tags           []*route53.Tag
}

func init() {
	register("Route53ResourceRecordSet", ListRoute53ResourceRecordSets)
}

func ListRoute53ResourceRecordSets(sess *session.Session) ([]Resource, error) {
	svc := route53.New(sess)

	resources := make([]Resource, 0)

	sub, err := ListRoute53HostedZones(sess)
	if err != nil {
		return nil, err
	}

	for _, resource := range sub {
		zone := resource.(*Route53HostedZone)
		rrs, err := ListResourceRecordsForZone(svc, zone.id, zone.name)
		if err != nil {
			return nil, err
		}

		resources = append(resources, rrs...)
	}

	return resources, nil
}

func ListResourceRecordsForZone(svc *route53.Route53, zoneId *string, zoneName *string) ([]Resource, error) {
	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId: zoneId,
	}

	hostedZoneTags, err := svc.ListTagsForResource(&route53.ListTagsForResourceInput{
		ResourceId:   zoneId,
		ResourceType: aws.String("hostedzone"),
	})

	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListResourceRecordSets(params)
		if err != nil {
			return nil, err
		}

		for _, rrs := range resp.ResourceRecordSets {
			resources = append(resources, &Route53ResourceRecordSet{
				svc:            svc,
				hostedZoneId:   zoneId,
				hostedZoneName: zoneName,
				data:           rrs,
				tags:           hostedZoneTags.ResourceTagSet.Tags,
			})
		}

		// make sure to list all with more than 100 records
		if *resp.IsTruncated {
			params.StartRecordName = resp.NextRecordName
			continue
		}

		break
	}

	return resources, nil
}

func (r *Route53ResourceRecordSet) Filter() error {
	if *r.data.Type == "NS" && *r.hostedZoneName == *r.data.Name {
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
				{
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

func (r *Route53ResourceRecordSet) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tag := range r.tags {
		properties.SetTagWithPrefix("hz", tag.Key, tag.Value)
	}
	properties.Set("Name", r.data.Name)
	properties.Set("Type", r.data.Type)
	return properties
}

func (r *Route53ResourceRecordSet) String() string {
	return *r.data.Name
}
