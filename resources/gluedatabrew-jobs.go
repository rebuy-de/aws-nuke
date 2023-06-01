package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/gluedatabrew"
)

type GlueDataBrewJobs struct {
	svc  *gluedatabrew.GlueDataBrew
	name *string
}

func init() {
	register("GlueDataBrewJobs", ListGlueDataBrewJobs)
}

func ListGlueDataBrewJobs(sess *session.Session) ([]Resource, error) {
	svc := gluedatabrew.New(sess)
	resources := []Resource{}

	params := &gluedatabrew.ListJobsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListJobs(params)
		if err != nil {
			return nil, err
		}

		for _, jobs := range output.Jobs {
			resources = append(resources, &GlueDataBrewJobs{
				svc:  svc,
				name: jobs.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueDataBrewJobs) Remove() error {
	_, err := f.svc.DeleteJob(&gluedatabrew.DeleteJobInput{
		Name: f.name,
	})

	return err
}

func (f *GlueDataBrewJobs) String() string {
	return *f.name
}
