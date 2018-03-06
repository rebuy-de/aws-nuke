package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudsearch"
)

type CloudSearchDomain struct {
	svc        *cloudsearch.CloudSearch
	domainName *string
}

func init() {
	register("CloudSearchDomain", ListCloudSearchDomains)
}

func ListCloudSearchDomains(sess *session.Session) ([]Resource, error) {
	svc := cloudsearch.New(sess)

	params := &cloudsearch.DescribeDomainsInput{}

	resp, err := svc.DescribeDomains(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, domain := range resp.DomainStatusList {
		resources = append(resources, &CloudSearchDomain{
			svc:        svc,
			domainName: domain.DomainName,
		})
	}
	return resources, nil
}

func (f *CloudSearchDomain) Remove() error {

	_, err := f.svc.DeleteDomain(&cloudsearch.DeleteDomainInput{
		DomainName: f.domainName,
	})

	return err
}

func (f *CloudSearchDomain) String() string {
	return *f.domainName
}
