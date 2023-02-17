package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opensearchservice"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type OSDomain struct {
	svc             *opensearchservice.OpenSearchService
	domainName      *string
	lastUpdatedTime *time.Time
	tagList         []*opensearchservice.Tag
}

func init() {
	register("OSDomain", ListOSDomains)
}

func ListOSDomains(sess *session.Session) ([]Resource, error) {
	svc := opensearchservice.New(sess)

	listResp, err := svc.ListDomainNames(&opensearchservice.ListDomainNamesInput{})
	if err != nil {
		return nil, err
	}
	var domainNames []*string
	for _, domain := range listResp.DomainNames {
		domainNames = append(domainNames, domain.DomainName)
	}

	resources := make([]Resource, 0)

	// early return to prevent the `missing required field, DescribeDomainsInput.DomainNames.` error
	if len(domainNames) == 0 {
		return resources, nil
	}

	descResp, err := svc.DescribeDomains(
		&opensearchservice.DescribeDomainsInput{
			DomainNames: domainNames,
		})
	if err != nil {
		return nil, err
	}

	for _, domain := range descResp.DomainStatusList {
		configResp, err := svc.DescribeDomainConfig(&opensearchservice.DescribeDomainConfigInput{DomainName: domain.DomainName})
		if err != nil {
			return nil, err
		}

		lto, err := svc.ListTags(&opensearchservice.ListTagsInput{ARN: domain.ARN})
		if err != nil {
			return nil, err
		}

		resources = append(resources, &OSDomain{
			svc:             svc,
			domainName:      domain.DomainName,
			lastUpdatedTime: configResp.DomainConfig.ClusterConfig.Status.UpdateDate,
			tagList:         lto.TagList,
		})
	}

	return resources, nil
}

func (o *OSDomain) Remove() error {
	_, err := o.svc.DeleteDomain(&opensearchservice.DeleteDomainInput{
		DomainName: o.domainName,
	})

	return err
}

func (o *OSDomain) Properties() types.Properties {
	properties := types.NewProperties().
		Set("LastUpdatedTime", o.lastUpdatedTime.Format(time.RFC3339))
	for _, t := range o.tagList {
		properties.SetTag(t.Key, t.Value)
	}
	return properties
}

func (o *OSDomain) String() string {
	return *o.domainName
}
