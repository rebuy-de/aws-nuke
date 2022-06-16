package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type IAMSigningCertificate struct {
	svc           *iam.IAM
	certificateId *string
	userName      *string
	status        *string
}

func init() {
	register("IAMSigningCertificate", ListIAMSigningCertificates)
}

func ListIAMSigningCertificates(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)
	resources := []Resource{}

	params := &iam.ListUsersInput{
		MaxItems: aws.Int64(100),
	}

	for {
		resp, err := svc.ListUsers(params)
		if err != nil {
			return nil, err
		}

		for _, out := range resp.Users {
			resp, err := svc.ListSigningCertificates(&iam.ListSigningCertificatesInput{
				UserName: out.UserName,
			})
			if err != nil {
				return nil, err
			}

			for _, signingCert := range resp.Certificates {
				resources = append(resources, &IAMSigningCertificate{
					svc:           svc,
					certificateId: signingCert.CertificateId,
					userName:      signingCert.UserName,
					status:        signingCert.Status,
				})
			}
		}

		if resp.Marker == nil {
			break
		}

		params.Marker = resp.Marker
	}

	return resources, nil
}

func (i *IAMSigningCertificate) Remove() error {
	_, err := i.svc.DeleteSigningCertificate(&iam.DeleteSigningCertificateInput{
		CertificateId: i.certificateId,
		UserName:      i.userName,
	})
	return err
}

func (i *IAMSigningCertificate) Properties() types.Properties {
	return types.NewProperties().
		Set("UserName", i.userName).
		Set("CertificateId", i.certificateId).
		Set("Status", i.status)
}

func (i *IAMSigningCertificate) String() string {
	return fmt.Sprintf("%s -> %s", *i.userName, *i.certificateId)
}
