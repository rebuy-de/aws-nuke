package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloud9"
)

type Cloud9Environment struct {
	svc           *cloud9.Cloud9
	environmentID *string
}

func init() {
	register("Cloud9Environment", ListCloud9Environments)
}

func ListCloud9Environments(sess *session.Session) ([]Resource, error) {
	svc := cloud9.New(sess)
	resources := []Resource{}

	params := &cloud9.ListEnvironmentsInput{
		MaxResults: aws.Int64(25),
	}

	for {
		resp, err := svc.ListEnvironments(params)
		if err != nil {
			return nil, err
		}

		for _, environmentID := range resp.EnvironmentIds {
			resources = append(resources, &Cloud9Environment{
				svc:           svc,
				environmentID: environmentID,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *Cloud9Environment) Remove() error {

	_, err := f.svc.DeleteEnvironment(&cloud9.DeleteEnvironmentInput{
		EnvironmentId: f.environmentID,
	})

	return err
}

func (f *Cloud9Environment) String() string {
	return *f.environmentID
}
