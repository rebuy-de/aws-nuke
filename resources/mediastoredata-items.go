package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediastore"
	"github.com/aws/aws-sdk-go/service/mediastoredata"
)

type MediaStoreDataItems struct {
	svc  *mediastoredata.MediaStoreData
	path *string
}

func init() {
	register("MediaStoreDataItems", ListMediaStoreDataItems)
}

func ListMediaStoreDataItems(sess *session.Session) ([]Resource, error) {
	containerSvc := mediastore.New(sess)
	svc := mediastoredata.New(sess)
	svc.ClientInfo.SigningName = "mediastore"
	resources := []Resource{}
	containers := []*mediastore.Container{}

	//List all containers
	containerParams := &mediastore.ListContainersInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := containerSvc.ListContainers(containerParams)
		if err != nil {
			return nil, err
		}

		for _, container := range output.Containers {
			containers = append(containers, container)
		}

		if output.NextToken == nil {
			break
		}

		containerParams.NextToken = output.NextToken
	}

	// List all Items per Container
	params := &mediastoredata.ListItemsInput{
		MaxResults: aws.Int64(100),
	}

	for _, container := range containers {
		if container.Endpoint == nil {
			continue
		}
		svc.Endpoint = *container.Endpoint
		output, err := svc.ListItems(params)
		if err != nil {
			return nil, err
		}

		for _, item := range output.Items {
			resources = append(resources, &MediaStoreDataItems{
				svc:  svc,
				path: item.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *MediaStoreDataItems) Remove() error {

	_, err := f.svc.DeleteObject(&mediastoredata.DeleteObjectInput{
		Path: f.path,
	})

	return err
}

func (f *MediaStoreDataItems) String() string {
	return *f.path
}
