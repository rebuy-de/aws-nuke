package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type SFNStateMachine struct {
	svc  *sfn.SFN
	ARN  *string
	name *string
	tags []*sfn.Tag
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
			tagsOutput, err := svc.ListTagsForResource(&sfn.ListTagsForResourceInput{
				ResourceArn: stateMachine.StateMachineArn,
			})

			if err != nil {
				return nil, err
			}

			resources = append(resources, &SFNStateMachine{
				svc:  svc,
				ARN:  stateMachine.StateMachineArn,
				name: stateMachine.Name,
				tags: tagsOutput.Tags,
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
	params := &sfn.ListExecutionsInput{
		StateMachineArn: f.ARN,
	}

	for {
		executions, execError := f.svc.ListExecutions(params)
		if execError != nil {
			break
		}
		for _, execs := range executions.Executions {

			f.svc.StopExecution(&sfn.StopExecutionInput{
				ExecutionArn: execs.ExecutionArn,
			})
		}

		if executions.NextToken == nil {
			break
		}
		params.NextToken = executions.NextToken
	}

	_, err := f.svc.DeleteStateMachine(&sfn.DeleteStateMachineInput{
		StateMachineArn: f.ARN,
	})

	return err
}

func (s *SFNStateMachine) Properties() types.Properties {
	properties := types.NewProperties()

	properties.Set("Name", s.name)

	for _, tagValue := range s.tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	return properties
}

func (f *SFNStateMachine) String() string {
	return *f.ARN
}
