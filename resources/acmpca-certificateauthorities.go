package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acmpca"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type ACMPCACertificateAuthority struct {
	svc    *acmpca.ACMPCA
	ARN    *string
	status *string
	tags   []*acmpca.Tag
}

func init() {
	register("ACMPCACertificateAuthority", ListACMPCACertificateAuthorities)
}

func ListACMPCACertificateAuthorities(sess *session.Session) ([]Resource, error) {
	svc := acmpca.New(sess)
	resources := []Resource{}
	tags := []*acmpca.Tag{}

	params := &acmpca.ListCertificateAuthoritiesInput{
		MaxResults: aws.Int64(100),
	}

	for {
		resp, err := svc.ListCertificateAuthorities(params)
		if err != nil {
			return nil, err
		}

		for _, certificateAuthority := range resp.CertificateAuthorities {
			tagParams := &acmpca.ListTagsInput{
				CertificateAuthorityArn: certificateAuthority.Arn,
				MaxResults:              aws.Int64(100),
			}

			for {
				tagResp, tagErr := svc.ListTags(tagParams)
				if tagErr != nil {
					return nil, tagErr
				}

				tags = append(tags, tagResp.Tags...)

				if tagResp.NextToken == nil {
					break
				}
				tagParams.NextToken = tagResp.NextToken
			}

			resources = append(resources, &ACMPCACertificateAuthority{
				svc:    svc,
				ARN:    certificateAuthority.Arn,
				status: certificateAuthority.Status,
				tags:   tags,
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

func (f *ACMPCACertificateAuthority) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tag := range f.tags {
		properties.SetTag(tag.Key, tag.Value)
	}
	properties.
		Set("ARN", f.ARN).
		Set("Status", f.status)
	return properties
}
