package resources

import "github.com/aws/aws-sdk-go/service/cloudformation"

func (n *CloudFormationNuke) ListStacks() ([]Resource, error) {
	resp, err := n.Service.DescribeStacks(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, stack := range resp.Stacks {
		resources = append(resources, &CloudFormationStack{
			svc:  n.Service,
			name: stack.StackName,
		})
	}
	return resources, nil
}

type CloudFormationStack struct {
	svc  *cloudformation.CloudFormation
	name *string
}

func (cfs *CloudFormationStack) Remove() error {
	_, err := cfs.svc.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: cfs.name,
	})
	return err
}

func (csf *CloudFormationStack) String() string {
	return *csf.name
}
