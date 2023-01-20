package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/guardduty"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type GuardDutyDetector struct {
	svc *guardduty.GuardDuty
	id  *string
}

func init() {
	register("GuardDutyDetector", ListGuardDutyDetectors)
}

func ListGuardDutyDetectors(sess *session.Session) ([]Resource, error) {
	svc := guardduty.New(sess)

	detectors := make([]Resource, 0)

	params := &guardduty.ListDetectorsInput{}

	err := svc.ListDetectorsPages(params, func(page *guardduty.ListDetectorsOutput, lastPage bool) bool {
		for _, out := range page.DetectorIds {
			detectors = append(detectors, &GuardDutyDetector{
				svc: svc,
				id:  out,
			})
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	return detectors, nil
}

func (detector *GuardDutyDetector) Remove() error {
	_, err := detector.svc.DeleteDetector(&guardduty.DeleteDetectorInput{
		DetectorId: detector.id,
	})
	return err
}

func (detector *GuardDutyDetector) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("DetectorID", detector.id)
	return properties
}

func (detector *GuardDutyDetector) String() string {
	return *detector.id
}
