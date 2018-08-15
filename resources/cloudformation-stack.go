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

	params := &cloudformation.DescribeStacksInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.DescribeStacks(params)
		if err != nil {
			return nil, err
		}
		for _, stack := range resp.Stacks {
			resources = append(resources, &CloudFormationStack{
				svc:   svc,
				stack: stack,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type CloudFormationStack struct {
	svc   *cloudformation.CloudFormation
	stack *cloudformation.Stack
}

func (cfs *CloudFormationStack) Remove() error {
	retainableResources, err := cfs.svc.ListStackResources(&cloudformation.ListStackResourcesInput{
		StackName: cfs.stack.StackName,
	})

	retain := make([]*string, 0)
	for _, r := range retainableResources.StackResourceSummaries {
		retain = append(retain, r.LogicalResourceId)
	}

	cfs.svc.DeleteStack(&cloudformation.DeleteStackInput{
		StackName:       cfs.stack.StackName,
		RetainResources: retain,
	})

	return err
}

func (cfs *CloudFormationStack) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", cfs.stack.StackName)
	for _, tagValue := range cfs.stack.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}

	return properties
}

func (cfs *CloudFormationStack) String() string {
	return *cfs.stack.StackName
}
