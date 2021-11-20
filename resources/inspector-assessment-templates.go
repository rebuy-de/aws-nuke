package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/inspector"
)

type InspectorAssessmentTemplate struct {
	svc *inspector.Inspector
	arn string
}

func init() {
	register("InspectorAssessmentTemplate", ListInspectorAssessmentTemplates)
}

func ListInspectorAssessmentTemplates(sess *session.Session) ([]Resource, error) {
	svc := inspector.New(sess)

	resp, err := svc.ListAssessmentTemplates(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.AssessmentTemplateArns {
		resources = append(resources, &InspectorAssessmentTemplate{
			svc: svc,
			arn: *out,
		})
	}

	return resources, nil
}

func (e *InspectorAssessmentTemplate) Remove() error {
	_, err := e.svc.DeleteAssessmentTemplate(&inspector.DeleteAssessmentTemplateInput{
		AssessmentTemplateArn: &e.arn,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *InspectorAssessmentTemplate) String() string {
	return e.arn
}
