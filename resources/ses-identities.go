package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SESIdentity struct {
	svc      *ses.SES
	identity *string
}

func init() {
	register("SESIdentity", ListSESIdentities)
}

func ListSESIdentities(sess *session.Session) ([]Resource, error) {
	svc := ses.New(sess)
	resources := []Resource{}

	params := &ses.ListIdentitiesInput{
		MaxItems: aws.Int64(100),
	}

	for {
		output, err := svc.ListIdentities(params)
		if err != nil {
			return nil, err
		}

		for _, identity := range output.Identities {
			resources = append(resources, &SESIdentity{
				svc:      svc,
				identity: identity,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *SESIdentity) Remove() error {

	_, err := f.svc.DeleteIdentity(&ses.DeleteIdentityInput{
		Identity: f.identity,
	})

	return err
}

func (f *SESIdentity) String() string {
	return *f.identity
}
