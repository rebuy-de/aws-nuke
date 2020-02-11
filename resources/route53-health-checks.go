package resources

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

func init() {
	register("Route53HealthCheck", ListRoute53HealthChecks)
}

func ListRoute53HealthChecks(sess *session.Session) ([]Resource, error) {
	svc := route53.New(sess)
	params := &route53.ListHealthChecksInput{}
	resources := make([]Resource, 0)

	var marker *string
	getHealthChecks := func() *bool {
		b := true
		return &b
	}()

	for *getHealthChecks == true {
		params.Marker = marker
		resp, err := svc.ListHealthChecks(params)
		if err != nil {
			return nil, err
		}

		for _, check := range resp.HealthChecks {
			resources = append(resources, &Route53HealthCheck{
				svc: svc,
				id:  check.Id,
			})
		}
		getHealthChecks = resp.IsTruncated
		marker = resp.NextMarker
	}

	return resources, nil
}

type Route53HealthCheck struct {
	svc *route53.Route53
	id  *string
}

func (hz *Route53HealthCheck) Remove() error {
	params := &route53.DeleteHealthCheckInput{
		HealthCheckId: hz.id,
	}

	_, err := hz.svc.DeleteHealthCheck(params)
	if err != nil {
		return err
	}

	return nil
}

func (hz *Route53HealthCheck) Properties() types.Properties {
	return types.NewProperties().
		Set("ID", hz.id)
}

func (hz *Route53HealthCheck) String() string {
	return fmt.Sprintf("%s", *hz.id)
}
