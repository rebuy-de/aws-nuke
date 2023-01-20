package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

func init() {
	register("ComprehendDocumentClassifier", ListComprehendDocumentClassifiers)
}

func ListComprehendDocumentClassifiers(sess *session.Session) ([]Resource, error) {
	svc := comprehend.New(sess)

	params := &comprehend.ListDocumentClassifiersInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListDocumentClassifiers(params)
		if err != nil {
			return nil, err
		}
		for _, documentClassifier := range resp.DocumentClassifierPropertiesList {
			resources = append(resources, &ComprehendDocumentClassifier{
				svc:                svc,
				documentClassifier: documentClassifier,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type ComprehendDocumentClassifier struct {
	svc                *comprehend.Comprehend
	documentClassifier *comprehend.DocumentClassifierProperties
}

func (ce *ComprehendDocumentClassifier) Remove() error {
	switch *ce.documentClassifier.Status {
	case "IN_ERROR":
		fallthrough
	case "TRAINED":
		{
			logrus.Infof("ComprehendDocumentClassifier deleteDocumentClassifier arn=%s status=%s", *ce.documentClassifier.DocumentClassifierArn, *ce.documentClassifier.Status)
			_, err := ce.svc.DeleteDocumentClassifier(&comprehend.DeleteDocumentClassifierInput{
				DocumentClassifierArn: ce.documentClassifier.DocumentClassifierArn,
			})
			return err
		}
	case "SUBMITTED":
		fallthrough
	case "TRAINING":
		{
			logrus.Infof("ComprehendDocumentClassifier stopTrainingDocumentClassifier arn=%s status=%s", *ce.documentClassifier.DocumentClassifierArn, *ce.documentClassifier.Status)
			_, err := ce.svc.StopTrainingDocumentClassifier(&comprehend.StopTrainingDocumentClassifierInput{
				DocumentClassifierArn: ce.documentClassifier.DocumentClassifierArn,
			})
			return err
		}
	default:
		{
			logrus.Infof("ComprehendDocumentClassifier already deleting arn=%s status=%s", *ce.documentClassifier.DocumentClassifierArn, *ce.documentClassifier.Status)
			return nil
		}
	}
}

func (ce *ComprehendDocumentClassifier) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("LanguageCode", ce.documentClassifier.LanguageCode)
	properties.Set("DocumentClassifierArn", ce.documentClassifier.DocumentClassifierArn)

	return properties
}

func (ce *ComprehendDocumentClassifier) String() string {
	return *ce.documentClassifier.DocumentClassifierArn
}
