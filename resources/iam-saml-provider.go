package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAMSAMLProvider struct {
	svc *iam.IAM
	arn string
}

func init() {
	register("IAMSAMLProvider", ListIAMSAMLProvider)
}

func ListIAMSAMLProvider(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)
	params := &iam.ListSAMLProvidersInput{}
	resources := make([]Resource, 0)

	resp, err := svc.ListSAMLProviders(params)
	if err != nil {
		return nil, err
	}

	for _, out := range resp.SAMLProviderList {
		resources = append(resources, &IAMSAMLProvider{
			svc: svc,
			arn: *out.Arn,
		})
	}

	return resources, nil
}

func (e *IAMSAMLProvider) Remove() error {
	_, err := e.svc.DeleteSAMLProvider(&iam.DeleteSAMLProviderInput{
		SAMLProviderArn: &e.arn,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMSAMLProvider) String() string {
	return e.arn
}
