package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
)

type GlueClassifier struct {
	svc  *glue.Glue
	name *string
}

func init() {
	register("GlueClassifier", ListGlueClassifiers)
}

func ListGlueClassifiers(sess *session.Session) ([]Resource, error) {
	svc := glue.New(sess)
	resources := []Resource{}

	params := &glue.GetClassifiersInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.GetClassifiers(params)
		if err != nil {
			return nil, err
		}

		for _, classifier := range output.Classifiers {
			switch {
			case classifier.GrokClassifier != nil:
				resources = append(resources, &GlueClassifier{
					svc:  svc,
					name: classifier.GrokClassifier.Name,
				})
			case classifier.JsonClassifier != nil:
				resources = append(resources, &GlueClassifier{
					svc:  svc,
					name: classifier.JsonClassifier.Name,
				})
			case classifier.XMLClassifier != nil:
				resources = append(resources, &GlueClassifier{
					svc:  svc,
					name: classifier.XMLClassifier.Name,
				})
			}
		}
		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueClassifier) Remove() error {

	_, err := f.svc.DeleteClassifier(&glue.DeleteClassifierInput{
		Name: f.name,
	})

	return err
}

func (f *GlueClassifier) String() string {
	return *f.name
}
