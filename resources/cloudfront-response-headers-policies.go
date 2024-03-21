package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CloudFrontResponseHeadersPolicy struct {
	svc  *cloudfront.CloudFront
	ID   *string
	name *string
}

func init() {
	register("CloudFrontResponseHeadersPolicy", ListCloudFrontResponseHeadersPolicies)
}

func ListCloudFrontResponseHeadersPolicies(sess *session.Session) ([]Resource, error) {
	svc := cloudfront.New(sess)
	resources := []Resource{}
	params := &cloudfront.ListResponseHeadersPoliciesInput{}

	for {
		resp, err := svc.ListResponseHeadersPolicies(params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.ResponseHeadersPolicyList.Items {
			resources = append(resources, &CloudFrontResponseHeadersPolicy{
				svc:  svc,
				ID:   item.ResponseHeadersPolicy.Id,
				name: item.ResponseHeadersPolicy.ResponseHeadersPolicyConfig.Name,
			})
		}

		if resp.ResponseHeadersPolicyList.NextMarker == nil {
			break
		}

		params.Marker = resp.ResponseHeadersPolicyList.NextMarker
	}

	return resources, nil
}

func (f *CloudFrontResponseHeadersPolicy) Filter() error {
	if strings.HasPrefix(*f.name, "Managed-") {
		return fmt.Errorf("Cannot delete default CloudFront Response headers policy")
	}
	return nil
}

func (f *CloudFrontResponseHeadersPolicy) Remove() error {
	resp, err := f.svc.GetResponseHeadersPolicy(&cloudfront.GetResponseHeadersPolicyInput{
		Id: f.ID,
	})
	if err != nil {
		return err
	}

	_, err = f.svc.DeleteResponseHeadersPolicy(&cloudfront.DeleteResponseHeadersPolicyInput{
		Id:      f.ID,
		IfMatch: resp.ETag,
	})

	return err
}

func (f *CloudFrontResponseHeadersPolicy) String() string {
	return *f.name
}

func (f *CloudFrontResponseHeadersPolicy) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ID", f.ID)
	properties.Set("Name", f.name)
	return properties
}
