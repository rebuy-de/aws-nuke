package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/inspector"
)

type InspectorAssessmentTarget struct {
	svc *inspector.Inspector
	arn string
}

func init() {
	register("InspectorAssessmentTarget", ListInspectorAssessmentTargets)
}

func ListInspectorAssessmentTargets(sess *session.Session) ([]Resource, error) {
	svc := inspector.New(sess)

	resp, err := svc.ListAssessmentTargets(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, out := range resp.AssessmentTargetArns {
		resources = append(resources, &InspectorAssessmentTarget{
			svc: svc,
			arn: *out,
		})
	}

	return resources, nil
}

func (e *InspectorAssessmentTarget) Remove() error {
	_, err := e.svc.DeleteAssessmentTarget(&inspector.DeleteAssessmentTargetInput{
		AssessmentTargetArn: &e.arn,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *InspectorAssessmentTarget) String() string {
	return e.arn
}
