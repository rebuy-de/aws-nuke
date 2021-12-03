package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type ESDomain struct {
	svc        *elasticsearchservice.ElasticsearchService
	domainName *string
	tags []*elasticsearchservice.Tag
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
			tags: domain.Tags,
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
func (e *ESDomain) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range e.tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	return properties
}


func (f *ESDomain) String() string {
	return *f.domainName
}
