package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMOpenIDConnectProvider struct {
	svc *iam.IAM
	arn string
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
		resources = append(resources, &IAMOpenIDConnectProvider{
			svc: svc,
			arn: *out.Arn,
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
