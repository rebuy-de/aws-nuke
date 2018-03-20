package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworks"
)

type OpsWorksInstance struct {
	svc *opsworks.OpsWorks
	ID  *string
}

func init() {
	register("OpsWorksInstance", ListOpsWorksInstances)
}

func ListOpsWorksInstances(sess *session.Session) ([]Resource, error) {
	svc := opsworks.New(sess)
	resources := []Resource{}

	stackParams := &opsworks.DescribeStacksInput{}

	resp, err := svc.DescribeStacks(stackParams)
	if err != nil {
		return nil, err
	}

	instanceParams := &opsworks.DescribeInstancesInput{}
	for _, stack := range resp.Stacks {
		instanceParams.StackId = stack.StackId
		output, err := svc.DescribeInstances(instanceParams)
		if err != nil {
			return nil, err
		}

		for _, instance := range output.Instances {
			resources = append(resources, &OpsWorksInstance{
				svc: svc,
				ID:  instance.InstanceId,
			})
		}
	}

	return resources, nil
}

func (f *OpsWorksInstance) Remove() error {

	_, err := f.svc.DeleteInstance(&opsworks.DeleteInstanceInput{
		InstanceId: f.ID,
	})

	return err
}

func (f *OpsWorksInstance) String() string {
	return *f.ID
}
