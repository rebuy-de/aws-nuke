package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type IAMOpenIDConnectProvider struct {
	svc  *iam.IAM
	arn  string
	tags []*iam.Tag
}

func init() {
	register("IAMOpenIDConnectProvider", ListIAMOpenIDConnectProvider)
}

func ListIAMOpenIDConnectProvider(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)
	params := &iam.ListOpenIDConnectProvidersInput{}
	resources := make([]Resource, 0)

	resp, err := svc.ListOpenIDConnectProviders(params)
	if err != nil {
		return nil, err
	}

	for _, out := range resp.OpenIDConnectProviderList {

		params := &iam.GetOpenIDConnectProviderInput{
			OpenIDConnectProviderArn: out.Arn,
		}
		resp, err := svc.GetOpenIDConnectProvider(params)

		if err != nil {
			return nil, err
		}

		resources = append(resources, &IAMOpenIDConnectProvider{
			svc:  svc,
			arn:  *out.Arn,
			tags: resp.Tags,
		})
	}

	return resources, nil
}

func (e *IAMOpenIDConnectProvider) Remove() error {
	_, err := e.svc.DeleteOpenIDConnectProvider(&iam.DeleteOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: &e.arn,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMOpenIDConnectProvider) String() string {
	return e.arn
}

func (e *IAMOpenIDConnectProvider) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Arn", e.arn)

	for _, tag := range e.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
