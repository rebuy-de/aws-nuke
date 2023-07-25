package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
)

type SFNStateMachine struct {
	svc *sfn.SFN
	ARN *string
}

func init() {
	register("SFNStateMachine", ListSFNStateMachines)
}

func ListSFNStateMachines(sess *session.Session) ([]Resource, error) {
	svc := sfn.New(sess)
	resources := []Resource{}

	params := &sfn.ListStateMachinesInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListStateMachines(params)
		if err != nil {
			return nil, err
		}

		for _, stateMachine := range output.StateMachines {
			resources = append(resources, &SFNStateMachine{
				svc: svc,
				ARN: stateMachine.StateMachineArn,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *SFNStateMachine) Remove() error {

	_, err := f.svc.DeleteStateMachine(&sfn.DeleteStateMachineInput{
		StateMachineArn: f.ARN,
	})

	return err
}

func (f *SFNStateMachine) String() string {
	return *f.ARN
}
