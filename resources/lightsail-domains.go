package resources

import (
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lightsail"
)

type LightsailDomain struct {
	svc        *lightsail.Lightsail
	domainName *string
}

func init() {
	register("LightsailDomain", ListLightsailDomains)
}

func ListLightsailDomains(sess *session.Session) ([]Resource, error) {
	svc := lightsail.New(sess)
	resources := []Resource{}

	if sess.Config.Region == nil || *sess.Config.Region != endpoints.UsEast1RegionID {
		// LightsailDomain only supports us-east-1
		return resources, nil
	}

	params := &lightsail.GetDomainsInput{}

	for {
		output, err := svc.GetDomains(params)
		if err != nil {
			return nil, err
		}

		for _, domain := range output.Domains {
			resources = append(resources, &LightsailDomain{
				svc:        svc,
				domainName: domain.Name,
			})
		}

		if output.NextPageToken == nil {
			break
		}

		params.PageToken = output.NextPageToken
	}

	return resources, nil
}

func (f *LightsailDomain) Remove() error {

	_, err := f.svc.DeleteDomain(&lightsail.DeleteDomainInput{
		DomainName: f.domainName,
	})

	return err
}

func (f *LightsailDomain) String() string {
	return *f.domainName
}
