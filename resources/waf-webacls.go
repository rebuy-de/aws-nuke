package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/waf"
)

type WAFWebACL struct {
	svc *waf.WAF
	ID  *string
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
			resources = append(resources, &WAFWebACL{
				svc: svc,
				ID:  webACL.WebACLId,
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

func (f *WAFWebACL) String() string {
	return *f.ID
}
