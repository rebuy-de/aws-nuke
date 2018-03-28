package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mobile"
)

type MobileProject struct {
	svc       *mobile.Mobile
	projectID *string
}

func init() {
	register("MobileProject", ListMobileProjects)
}

func ListMobileProjects(sess *session.Session) ([]Resource, error) {
	svc := mobile.New(sess)
	svc.ClientInfo.SigningName = "AWSMobileHubService"
	resources := []Resource{}

	params := &mobile.ListProjectsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListProjects(params)
		if err != nil {
			return nil, err
		}

		for _, project := range output.Projects {
			resources = append(resources, &MobileProject{
				svc:       svc,
				projectID: project.ProjectId,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *MobileProject) Remove() error {

	_, err := f.svc.DeleteProject(&mobile.DeleteProjectInput{
		ProjectId: f.projectID,
	})

	return err
}

func (f *MobileProject) String() string {
	return *f.projectID
}
