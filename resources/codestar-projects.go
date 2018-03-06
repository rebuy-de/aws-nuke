package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codestar"
)

type CodeStarProject struct {
	svc *codestar.CodeStar
	id  *string
}

func init() {
	register("CodeStarProject", ListCodeStarProjects)
}

func ListCodeStarProjects(sess *session.Session) ([]Resource, error) {
	svc := codestar.New(sess)
	resources := []Resource{}

	params := &codestar.ListProjectsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListProjects(params)
		if err != nil {
			return nil, err
		}

		for _, project := range output.Projects {
			resources = append(resources, &CodeStarProject{
				svc: svc,
				id:  project.ProjectId,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *CodeStarProject) Remove() error {

	_, err := f.svc.DeleteProject(&codestar.DeleteProjectInput{
		Id: f.id,
	})

	return err
}

func (f *CodeStarProject) String() string {
	return *f.id
}
