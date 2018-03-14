package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SSMDocument struct {
	svc  *ssm.SSM
	name *string
}

func init() {
	register("SSMDocument", ListSSMDocuments)
}

func ListSSMDocuments(sess *session.Session) ([]Resource, error) {
	svc := ssm.New(sess)
	resources := []Resource{}

	documentKeyFilter := []*ssm.DocumentKeyValuesFilter{
		{
			Key:    aws.String("Owner"),
			Values: []*string{aws.String("Self")},
		},
	}

	params := &ssm.ListDocumentsInput{
		MaxResults: aws.Int64(50),
		Filters:    documentKeyFilter,
	}

	for {
		output, err := svc.ListDocuments(params)
		if err != nil {
			return nil, err
		}

		for _, documentIdentifier := range output.DocumentIdentifiers {
			resources = append(resources, &SSMDocument{
				svc:  svc,
				name: documentIdentifier.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *SSMDocument) Remove() error {

	_, err := f.svc.DeleteDocument(&ssm.DeleteDocumentInput{
		Name: f.name,
	})

	return err
}

func (f *SSMDocument) String() string {
	return *f.name
}
