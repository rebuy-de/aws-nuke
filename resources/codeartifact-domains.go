package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codeartifact"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CodeArtifactDomain struct {
	svc  *codeartifact.CodeArtifact
	name *string
	tags map[string]*string
}

func init() {
	register("CodeArtifactDomain", ListCodeArtifactDomains)
}

func ListCodeArtifactDomains(sess *session.Session) ([]Resource, error) {
	svc := codeartifact.New(sess)
	resources := []Resource{}

	params := &codeartifact.ListDomainsInput{}

	for {
		resp, err := svc.ListDomains(params)
		if err != nil {
			return nil, err
		}

		for _, domain := range resp.Domains {
			desc, err := svc.DescribeDomain(&codeartifact.DescribeDomainInput{Domain: domain.Name})
			if err != nil {
				return nil, err
			}

			resources = append(resources, &CodeArtifactDomain{
				svc:  svc,
				name: domain.Name,
				tags: GetDomainTags(svc, desc.Domain.Arn),
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func GetDomainTags(svc *codeartifact.CodeArtifact, arn *string) map[string]*string {
	tags := map[string]*string{}

	resp, _ := svc.ListTagsForResource(&codeartifact.ListTagsForResourceInput{ResourceArn: arn})
	for _, tag := range resp.Tags {
		tags[*tag.Key] = tag.Value
	}

	return tags
}

func (d *CodeArtifactDomain) Remove() error {
	_, err := d.svc.DeleteDomain(&codeartifact.DeleteDomainInput{
		Domain: d.name,
	})
	return err
}

func (d *CodeArtifactDomain) String() string {
	return *d.name
}

func (d *CodeArtifactDomain) Properties() types.Properties {
	properties := types.NewProperties()
	for key, tag := range d.tags {
		properties.SetTag(&key, tag)
	}
	properties.Set("Name", d.name)
	return properties
}
