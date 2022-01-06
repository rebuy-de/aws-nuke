package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/inspector"
)

type InspectorAssessmentRun struct {
	svc *inspector.Inspector
	arn string
}

func init() {
	register("InspectorAssessmentRun", ListInspectorAssessmentRuns)
}

func ListInspectorAssessmentRuns(sess *session.Session) ([]Resource, error) {
	svc := inspector.New(sess)

	resp, err := svc.ListAssessmentRuns(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.AssessmentRunArns {
		resources = append(resources, &InspectorAssessmentRun{
			svc: svc,
			arn: *out,
		})
	}

	return resources, nil
}

func (e *InspectorAssessmentRun) Remove() error {
	_, err := e.svc.DeleteAssessmentRun(&inspector.DeleteAssessmentRunInput{
		AssessmentRunArn: &e.arn,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *InspectorAssessmentRun) String() string {
	return e.arn
}
