package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acmpca"
)

type ACMPCACertificateAuthority struct {
	svc    *acmpca.ACMPCA
	ARN    *string
	status *string
}

func init() {
	register("ACMPCACertificateAuthority", ListACMPCACertificateAuthorities)
}

func ListACMPCACertificateAuthorities(sess *session.Session) ([]Resource, error) {
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
			resources = append(resources, &ACMPCACertificateAuthority{
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

func (f *ACMPCACertificateAuthority) Remove() error {

	_, err := f.svc.DeleteCertificateAuthority(&acmpca.DeleteCertificateAuthorityInput{
		CertificateAuthorityArn: f.ARN,
	})

	return err
}

func (f *ACMPCACertificateAuthority) String() string {
	return *f.ARN
}

func (f *ACMPCACertificateAuthority) Filter() error {
	if *f.status == "DELETED" {
		return fmt.Errorf("already deleted")
	} else {
		return nil
	}
}
