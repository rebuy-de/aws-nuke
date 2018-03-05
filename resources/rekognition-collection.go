package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

type RekognitionCollection struct {
	svc *rekognition.Rekognition
	id  *string
}

func init() {
	register("RekognitionCollection", ListRekognitionCollections)
}

func ListRekognitionCollections(sess *session.Session) ([]Resource, error) {
	svc := rekognition.New(sess)
	resources := []Resource{}

	params := &rekognition.ListCollectionsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListCollections(params)
		if err != nil {
			return nil, err
		}

		for _, collection := range output.CollectionIds {
			resources = append(resources, &RekognitionCollection{
				svc: svc,
				id:  collection,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *RekognitionCollection) Remove() error {

	_, err := f.svc.DeleteCollection(&rekognition.DeleteCollectionInput{
		CollectionId: f.id,
	})

	return err
}

func (f *RekognitionCollection) String() string {
	return *f.id
}
