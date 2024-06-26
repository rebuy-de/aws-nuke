package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFWebACL struct {
	svc  *waf.WAF
	ID   *string
	tags []*waf.Tag
}

func init() {
	register("WAFWebACL", ListWAFWebACLs)
}

func ListWAFWebACLs(sess *session.Session) ([]Resource, error) {
	svc := waf.New(sess)
	resources := []Resource{}

	params := &waf.ListWebACLsInput{
		Limit: aws.Int64(50),
	}

	for {
		resp, err := svc.ListWebACLs(params)
		if err != nil {
			return nil, err
		}

		for _, webACL := range resp.WebACLs {
			acl, err := svc.GetWebACL(&waf.GetWebACLInput{
				WebACLId: webACL.WebACLId,
			})
			if err != nil {
				return nil, err
			}
			tags, err := svc.ListTagsForResource(&waf.ListTagsForResourceInput{
				ResourceARN: acl.WebACL.WebACLArn,
			})
			if err != nil {
				return nil, err
			}

			resources = append(resources, &WAFWebACL{
				svc:  svc,
				ID:   webACL.WebACLId,
				tags: tags.TagInfoForResource.TagList,
			})
		}

		if resp.NextMarker == nil {
			break
		}

		params.NextMarker = resp.NextMarker
	}

	return resources, nil
}

func (f *WAFWebACL) Remove() error {

	tokenOutput, err := f.svc.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	_, err = f.svc.DeleteWebACL(&waf.DeleteWebACLInput{
		WebACLId:    f.ID,
		ChangeToken: tokenOutput.ChangeToken,
	})

	return err
}

func (f *WAFWebACL) Properties() types.Properties {
	properties := types.NewProperties()

	for _, tag := range f.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}

func (f *WAFWebACL) String() string {
	return *f.ID
}
