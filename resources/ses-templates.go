package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SESTemplate struct {
	svc  *ses.SES
	name *string
}

func init() {
	register("SESTemplate", ListSESTemplates)
}

func ListSESTemplates(sess *session.Session) ([]Resource, error) {
	svc := ses.New(sess)
	resources := []Resource{}

	params := &ses.ListTemplatesInput{
		MaxItems: aws.Int64(100),
	}

	for {
		output, err := svc.ListTemplates(params)
		if err != nil {
			return nil, err
		}

		for _, templateMetadata := range output.TemplatesMetadata {
			resources = append(resources, &SESTemplate{
				svc:  svc,
				name: templateMetadata.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *SESTemplate) Remove() error {

	_, err := f.svc.DeleteTemplate(&ses.DeleteTemplateInput{
		TemplateName: f.name,
	})

	return err
}

func (f *SESTemplate) String() string {
	return *f.name
}
