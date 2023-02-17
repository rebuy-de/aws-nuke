package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/gluedatabrew"
)

type GlueDataBrewProjects struct {
	svc  *gluedatabrew.GlueDataBrew
	name *string
}

func init() {
	register("GlueDataBrewProjects", ListGlueDataBrewProjects)
}

func ListGlueDataBrewProjects(sess *session.Session) ([]Resource, error) {
	svc := gluedatabrew.New(sess)
	resources := []Resource{}

	params := &gluedatabrew.ListProjectsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListProjects(params)
		if err != nil {
			return nil, err
		}

		for _, project := range output.Projects {
			resources = append(resources, &GlueDataBrewProjects{
				svc:  svc,
				name: project.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueDataBrewProjects) Remove() error {
	_, err := f.svc.DeleteProject(&gluedatabrew.DeleteProjectInput{
		Name: f.name,
	})

	return err
}

func (f *GlueDataBrewProjects) String() string {
	return *f.name
}
