package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
)

type ESDomain struct {
	svc        *elasticsearchservice.ElasticsearchService
	domainName *string
	tagList    []*elasticsearchservice.Tag
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
		dedo, err := svc.DescribeElasticsearchDomain(
			&elasticsearchservice.DescribeElasticsearchDomainInput{DomainName: domain.DomainName})
		if err != nil {
			return nil, err
		}
		lto, err := svc.ListTags(&elasticsearchservice.ListTagsInput{ARN: dedo.DomainStatus.ARN})
		if err != nil {
			return nil, err
		}
		resources = append(resources, &ESDomain{
			svc:        svc,
			domainName: domain.DomainName,
			tagList:    lto.TagList,
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
