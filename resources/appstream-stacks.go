package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appstream"
)

type AppStreamStack struct {
	svc  *appstream.AppStream
	name *string
}

func init() {
	register("AppStreamStack", ListAppStreamStacks)
}

func ListAppStreamStacks(sess *session.Session) ([]Resource, error) {
	svc := appstream.New(sess)
	resources := []Resource{}

	params := &appstream.DescribeStacksInput{}

	for {
		output, err := svc.DescribeStacks(params)
		if err != nil {
			return nil, err
		}

		for _, stack := range output.Stacks {
			resources = append(resources, &AppStreamStack{
				svc:  svc,
				name: stack.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *AppStreamStack) Remove() error {

	_, err := f.svc.DeleteStack(&appstream.DeleteStackInput{
		Name: f.name,
	})

	return err
}

func (f *AppStreamStack) String() string {
	return *f.name
}
