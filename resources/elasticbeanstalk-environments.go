package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ElasticBeanstalkEnvironment struct {
	svc  *elasticbeanstalk.ElasticBeanstalk
	ID   *string
	name *string
}

func init() {
	register("ElasticBeanstalkEnvironment", ListElasticBeanstalkEnvironments)
}

func ListElasticBeanstalkEnvironments(sess *session.Session) ([]Resource, error) {
	svc := elasticbeanstalk.New(sess)
	resources := []Resource{}

	params := &elasticbeanstalk.DescribeEnvironmentsInput{
		MaxRecords:     aws.Int64(100),
		IncludeDeleted: aws.Bool(false),
	}

	for {
		output, err := svc.DescribeEnvironments(params)
		if err != nil {
			return nil, err
		}

		for _, environment := range output.Environments {
			resources = append(resources, &ElasticBeanstalkEnvironment{
				svc:  svc,
				ID:   environment.EnvironmentId,
				name: environment.EnvironmentName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *ElasticBeanstalkEnvironment) Remove() error {

	_, err := f.svc.TerminateEnvironment(&elasticbeanstalk.TerminateEnvironmentInput{
		EnvironmentId:      f.ID,
		ForceTerminate:     aws.Bool(true),
		TerminateResources: aws.Bool(true),
	})

	return err
}

func (e *ElasticBeanstalkEnvironment) Properties() types.Properties {
	return types.NewProperties().
		Set("Name", e.name)
}

func (f *ElasticBeanstalkEnvironment) String() string {
	return *f.ID
}
