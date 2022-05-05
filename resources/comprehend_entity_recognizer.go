package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

func init() {
	register("ComprehendEntityRecognizer", ListComprehendEntityRecognizers)
}

func ListComprehendEntityRecognizers(sess *session.Session) ([]Resource, error) {
	svc := comprehend.New(sess)

	params := &comprehend.ListEntityRecognizersInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListEntityRecognizers(params)
		if err != nil {
			return nil, err
		}
		for _, entityRecognizer := range resp.EntityRecognizerPropertiesList {
			resources = append(resources, &ComprehendEntityRecognizer{
				svc:              svc,
				entityRecognizer: entityRecognizer,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type ComprehendEntityRecognizer struct {
	svc              *comprehend.Comprehend
	entityRecognizer *comprehend.EntityRecognizerProperties
}

func (ce *ComprehendEntityRecognizer) Remove() error {
	switch *ce.entityRecognizer.Status {
	case "IN_ERROR":
		fallthrough
	case "TRAINED":
		{
			logrus.Infof("ComprehendEntityRecognizer deleteEntityRecognizer arn=%s status=%s", *ce.entityRecognizer.EntityRecognizerArn, *ce.entityRecognizer.Status)
			_, err := ce.svc.DeleteEntityRecognizer(&comprehend.DeleteEntityRecognizerInput{
				EntityRecognizerArn: ce.entityRecognizer.EntityRecognizerArn,
			})
			return err
		}
	case "SUBMITTED":
		fallthrough
	case "TRAINING":
		{
			logrus.Infof("ComprehendEntityRecognizer stopTrainingEntityRecognizer arn=%s status=%s", *ce.entityRecognizer.EntityRecognizerArn, *ce.entityRecognizer.Status)
			_, err := ce.svc.StopTrainingEntityRecognizer(&comprehend.StopTrainingEntityRecognizerInput{
				EntityRecognizerArn: ce.entityRecognizer.EntityRecognizerArn,
			})
			return err
		}
	default:
		{
			logrus.Infof("ComprehendEntityRecognizer already deleting arn=%s status=%s", *ce.entityRecognizer.EntityRecognizerArn, *ce.entityRecognizer.Status)
			return nil
		}
	}
}

func (ce *ComprehendEntityRecognizer) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("LanguageCode", ce.entityRecognizer.LanguageCode)
	properties.Set("EntityRecognizerArn", ce.entityRecognizer.EntityRecognizerArn)

	return properties
}

func (ce *ComprehendEntityRecognizer) String() string {
	return *ce.entityRecognizer.EntityRecognizerArn
}
