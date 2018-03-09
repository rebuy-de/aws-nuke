package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
)

type ElasticBeanstalkApplication struct {
	svc  *elasticbeanstalk.ElasticBeanstalk
	name *string
}

func init() {
	register("ElasticBeanstalkApplication", ListElasticBeanstalkApplications)
}

func ListElasticBeanstalkApplications(sess *session.Session) ([]Resource, error) {
	svc := elasticbeanstalk.New(sess)
	resources := []Resource{}

	params := &elasticbeanstalk.DescribeApplicationsInput{}

	output, err := svc.DescribeApplications(params)
	if err != nil {
		return nil, err
	}

	for _, application := range output.Applications {
		resources = append(resources, &ElasticBeanstalkApplication{
			svc:  svc,
			name: application.ApplicationName,
		})
	}

	return resources, nil
}

func (f *ElasticBeanstalkApplication) Remove() error {

	_, err := f.svc.DeleteApplication(&elasticbeanstalk.DeleteApplicationInput{
		ApplicationName: f.name,
	})

	return err
}

func (f *ElasticBeanstalkApplication) String() string {
	return *f.name
}
