package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type ESDomain struct {
	svc        *elasticsearchservice.ElasticsearchService
	domainName *string
	tags map[string]*string

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
		tags, err := svc.ListTags(&elasticsearchservice.ListTagsInput{
			Resource: domain.FunctionArn,
		})

		if err != nil {
			continue
		}
		resources = append(resources, &ESDomain{
			svc:        svc,
			domainName: domain.DomainName,
			tags: tags.Tags,
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
func (f *ESDomain) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("DomainName", f.domainName)

	for key, val := range f.tags {
		properties.SetTag(&key, val)
	}

	return properties
}


func (f *ESDomain) String() string {
	return *f.domainName
}
