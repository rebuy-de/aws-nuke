package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acmpca"
)

type ACMPCACertificateAuthorityState struct {
	svc    *acmpca.ACMPCA
	ARN    *string
	status *string
}

func init() {
	register("ACMPCACertificateAuthorityState", ListACMPCACertificateAuthorityStates)
}

func ListACMPCACertificateAuthorityStates(sess *session.Session) ([]Resource, error) {
	svc := acmpca.New(sess)
	resources := []Resource{}

	params := &acmpca.ListCertificateAuthoritiesInput{
		MaxResults: aws.Int64(100),
	}

	for {
		resp, err := svc.ListCertificateAuthorities(params)
		if err != nil {
			return nil, err
		}

		for _, certificateAuthority := range resp.CertificateAuthorities {
			resources = append(resources, &ACMPCACertificateAuthorityState{
				svc:    svc,
				ARN:    certificateAuthority.Arn,
				status: certificateAuthority.Status,
			})
		}
		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}
	return resources, nil
}

func (f *ACMPCACertificateAuthorityState) Remove() error {

	_, err := f.svc.UpdateCertificateAuthority(&acmpca.UpdateCertificateAuthorityInput{
		CertificateAuthorityArn: f.ARN,
		Status:                  aws.String("DISABLED"),
	})

	return err
}

func (f *ACMPCACertificateAuthorityState) String() string {
	return *f.ARN
}

func (f *ACMPCACertificateAuthorityState) Filter() error {

	switch *f.status {
	case "CREATING":
		return fmt.Errorf("available for deletion")
	case "PENDING_CERTIFICATE":
		return fmt.Errorf("available for deletion")
	case "DISABLED":
		return fmt.Errorf("available for deletion")
	case "DELETED":
		return fmt.Errorf("already deleted")
	default:
		return nil
	}

}
