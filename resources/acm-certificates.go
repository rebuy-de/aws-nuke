package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
)

type ACMCertificate struct {
	svc            *acm.ACM
	certificateARN *string
}

func init() {
	register("ACMCertificate", ListACMCertificates)
}

func ListACMCertificates(sess *session.Session) ([]Resource, error) {
	svc := acm.New(sess)
	resources := []Resource{}

	params := &acm.ListCertificatesInput{
		MaxItems: aws.Int64(100),
	}

	for {
		resp, err := svc.ListCertificates(params)
		if err != nil {
			return nil, err
		}

		for _, certificate := range resp.CertificateSummaryList {
			resources = append(resources, &ACMCertificate{
				svc:            svc,
				certificateARN: certificate.CertificateArn,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *ACMCertificate) Remove() error {

	_, err := f.svc.DeleteCertificate(&acm.DeleteCertificateInput{
		CertificateArn: f.certificateARN,
	})

	return err
}

func (f *ACMCertificate) String() string {
	return *f.certificateARN
}
