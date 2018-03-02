package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
)

type ESDomain struct {
	svc        *elasticsearchservice.ElasticsearchService
	domainName *string
}

func init() {
	register("ESDomain", ListESDomains)
}

func ListESDomains(sess *session.Session) ([]Resource, error) {
	svc := elasticsearchservice.New(sess)

	params := &elasticsearchservice.ListDomainNamesInput{}
	resp, err := svc.ListDomainNames(params)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, domain := range resp.DomainNames {
		resources = append(resources, &ESDomain{
			svc:        svc,
			domainName: domain.DomainName,
		})
	}

	return resources, nil
}

func (f *ESDomain) Remove() error {

	_, err := f.svc.DeleteElasticsearchDomain(&elasticsearchservice.DeleteElasticsearchDomainInput{
		DomainName: f.domainName,
	})

	return err
}

func (f *ESDomain) String() string {
	return *f.domainName
}
