package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

type ELBv2TargetGroup struct {
	svc  *elbv2.ELBV2
	name *string
	arn  *string
}

func init() {
	register("ELBv2TargetGroup", ListELBv2TargetGroups)
}

func ListELBv2TargetGroups(sess *session.Session) ([]Resource, error) {
	svc := elbv2.New(sess)

	resp, err := svc.DescribeTargetGroups(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, elbv2TargetGroup := range resp.TargetGroups {
		resources = append(resources, &ELBv2TargetGroup{
			svc:  svc,
			name: elbv2TargetGroup.TargetGroupName,
			arn:  elbv2TargetGroup.TargetGroupArn,
		})
	}

	return resources, nil
}

func (e *ELBv2TargetGroup) Remove() error {
	_, err := e.svc.DeleteTargetGroup(&elbv2.DeleteTargetGroupInput{
		TargetGroupArn: e.arn,
	})

	if err != nil {
		return err
	}

	return nil
}

func (e *ELBv2TargetGroup) String() string {
	return *e.name
}
