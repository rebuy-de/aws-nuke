package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type SageMakerSpace struct {
	svc              *sagemaker.SageMaker
	domainID         *string
	spaceDisplayName *string
	spaceName        *string
	status           *string
	lastModifiedTime *time.Time
}

func init() {
	register("SageMakerSpace", ListSageMakerSpaces)
}

func ListSageMakerSpaces(sess *session.Session) ([]Resource, error) {
	svc := sagemaker.New(sess)
	resources := []Resource{}

	params := &sagemaker.ListSpacesInput{
		MaxResults: aws.Int64(30),
	}

	for {
		resp, err := svc.ListSpaces(params)
		if err != nil {
			return nil, err
		}

		for _, space := range resp.Spaces {
			resources = append(resources, &SageMakerSpace{
				svc:              svc,
				domainID:         space.DomainId,
				spaceDisplayName: space.SpaceDisplayName,
				spaceName:        space.SpaceName,
				status:           space.Status,
				lastModifiedTime: space.LastModifiedTime,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *SageMakerSpace) Remove() error {
	_, err := f.svc.DeleteSpace(&sagemaker.DeleteSpaceInput{
		DomainId:  f.domainID,
		SpaceName: f.spaceName,
	})

	return err
}

func (f *SageMakerSpace) String() string {
	return *f.spaceName
}

func (i *SageMakerSpace) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("DomainID", i.domainID).
		Set("SpaceDisplayName", i.spaceDisplayName).
		Set("SpaceName", i.spaceName).
		Set("Status", i.status).
		Set("LastModifiedTime", i.lastModifiedTime)
	return properties
}
