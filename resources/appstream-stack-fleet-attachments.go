package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appstream"
)

type AppStreamStackFleetAttachment struct {
	svc       *appstream.AppStream
	stackName *string
	fleetName *string
}

func init() {
	register("AppStreamStackFleetAttachment", ListAppStreamStackFleetAttachments)
}

func ListAppStreamStackFleetAttachments(sess *session.Session) ([]Resource, error) {
	svc := appstream.New(sess)
	resources := []Resource{}
	stacks := []*appstream.Stack{}
	params := &appstream.DescribeStacksInput{}

	for {
		output, err := svc.DescribeStacks(params)
		if err != nil {
			return nil, err
		}

		for _, stack := range output.Stacks {
			stacks = append(stacks, stack)
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	stackAssocParams := &appstream.ListAssociatedFleetsInput{}
	for _, stack := range stacks {

		stackAssocParams.StackName = stack.Name
		output, err := svc.ListAssociatedFleets(stackAssocParams)
		if err != nil {
			return nil, err
		}

		for _, name := range output.Names {
			resources = append(resources, &AppStreamStackFleetAttachment{
				svc:       svc,
				stackName: stack.Name,
				fleetName: name,
			})
		}
	}

	return resources, nil
}

func (f *AppStreamStackFleetAttachment) Remove() error {

	_, err := f.svc.DisassociateFleet(&appstream.DisassociateFleetInput{
		StackName: f.stackName,
		FleetName: f.fleetName,
	})

	return err
}

func (f *AppStreamStackFleetAttachment) String() string {
	return fmt.Sprintf("%s -> %s", *f.stackName, *f.fleetName)
}
