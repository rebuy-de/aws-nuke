package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SESConfigurationSet struct {
	svc  *ses.SES
	name *string
}

func init() {
	register("SESConfigurationSet", ListSESConfigurationSets)
}

func ListSESConfigurationSets(sess *session.Session) ([]Resource, error) {
	svc := ses.New(sess)
	resources := []Resource{}

	params := &ses.ListConfigurationSetsInput{
		MaxItems: aws.Int64(100),
	}

	for {
		output, err := svc.ListConfigurationSets(params)
		if err != nil {
			return nil, err
		}

		for _, configurationSet := range output.ConfigurationSets {
			resources = append(resources, &SESConfigurationSet{
				svc:  svc,
				name: configurationSet.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *SESConfigurationSet) Remove() error {

	_, err := f.svc.DeleteConfigurationSet(&ses.DeleteConfigurationSetInput{
		ConfigurationSetName: f.name,
	})

	return err
}

func (f *SESConfigurationSet) String() string {
	return *f.name
}
