package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
)

type GlueTrigger struct {
	svc  *glue.Glue
	name *string
}

func init() {
	register("GlueTrigger", ListGlueTriggers)
}

func ListGlueTriggers(sess *session.Session) ([]Resource, error) {
	svc := glue.New(sess)
	resources := []Resource{}

	params := &glue.GetTriggersInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.GetTriggers(params)
		if err != nil {
			return nil, err
		}

		for _, trigger := range output.Triggers {
			resources = append(resources, &GlueTrigger{
				svc:  svc,
				name: trigger.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueTrigger) Remove() error {

	_, err := f.svc.DeleteTrigger(&glue.DeleteTriggerInput{
		Name: f.name,
	})

	return err
}

func (f *GlueTrigger) String() string {
	return *f.name
}
