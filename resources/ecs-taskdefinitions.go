package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type ECSTaskDefinition struct {
	svc *ecs.ECS
	ARN *string
}

func init() {
	register("ECSTaskDefinition", ListECSTaskDefinitions)
}

func ListECSTaskDefinitions(sess *session.Session) ([]Resource, error) {
	svc := ecs.New(sess)
	resources := []Resource{}

	params := &ecs.ListTaskDefinitionsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListTaskDefinitions(params)
		if err != nil {
			return nil, err
		}

		for _, taskDefinitionARN := range output.TaskDefinitionArns {
			resources = append(resources, &ECSTaskDefinition{
				svc: svc,
				ARN: taskDefinitionARN,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *ECSTaskDefinition) Remove() error {

	_, err := f.svc.DeregisterTaskDefinition(&ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: f.ARN,
	})
	if err != nil {
		return err
	}
	taskDefinitions := make([]*string, 0)
	taskDefinitions = append(taskDefinitions, f.ARN)
	_, err1 := f.svc.DeleteTaskDefinitions(&ecs.DeleteTaskDefinitionsInput{
		TaskDefinitions: taskDefinitions,
	})
	if err1 != nil {
		return err1
	}
	return nil
}

func (f *ECSTaskDefinition) String() string {
	return *f.ARN
}
