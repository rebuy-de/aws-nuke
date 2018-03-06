package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/resourcegroups"
)

type ResourceGroupGroup struct {
	svc       *resourcegroups.ResourceGroups
	groupName *string
}

func init() {
	register("ResourceGroupGroup", ListResourceGroupGroups)
}

func ListResourceGroupGroups(sess *session.Session) ([]Resource, error) {
	svc := resourcegroups.New(sess)
	resources := []Resource{}

	params := &resourcegroups.ListGroupsInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.ListGroups(params)
		if err != nil {
			return nil, err
		}

		for _, group := range output.Groups {
			resources = append(resources, &ResourceGroupGroup{
				svc:       svc,
				groupName: group.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *ResourceGroupGroup) Remove() error {

	_, err := f.svc.DeleteGroup(&resourcegroups.DeleteGroupInput{
		GroupName: f.groupName,
	})

	return err
}

func (f *ResourceGroupGroup) String() string {
	return *f.groupName
}
