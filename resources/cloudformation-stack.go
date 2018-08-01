package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

func init() {
	register("CloudFormationStack", ListCloudFormationStacks)
}

func ListCloudFormationStacks(sess *session.Session) ([]Resource, error) {
	svc := cloudformation.New(sess)

	resp, err := svc.DescribeStacks(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, stack := range resp.Stacks {
		resources = append(resources, &CloudFormationStack{
			svc:   svc,
			stack: stack,
		})
	}
	return resources, nil
}

type CloudFormationStack struct {
	svc   *cloudformation.CloudFormation
	stack *cloudformation.Stack
}

func (cfs *CloudFormationStack) Remove() error {
	_, err := cfs.svc.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: cfs.stack.StackName,
	})
	return err
}

func (cfs *CloudFormationStack) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("Name", cfs.stack.StackName)
	for _, tagValue := range cfs.stack.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}

	return properties
}

func (cfs *CloudFormationStack) String() string {
	return *cfs.stack.StackName
}
