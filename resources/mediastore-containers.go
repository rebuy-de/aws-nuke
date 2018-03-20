package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediastore"
)

type MediaStoreContainer struct {
	svc  *mediastore.MediaStore
	name *string
}

func init() {
	register("MediaStoreContainer", ListMediaStoreContainers)
}

func ListMediaStoreContainers(sess *session.Session) ([]Resource, error) {
	svc := mediastore.New(sess)
	resources := []Resource{}

	params := &mediastore.ListContainersInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListContainers(params)
		if err != nil {
			return nil, err
		}

		for _, container := range output.Containers {
			resources = append(resources, &MediaStoreContainer{
				svc:  svc,
				name: container.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *MediaStoreContainer) Remove() error {

	_, err := f.svc.DeleteContainer(&mediastore.DeleteContainerInput{
		ContainerName: f.name,
	})

	return err
}

func (f *MediaStoreContainer) String() string {
	return *f.name
}
