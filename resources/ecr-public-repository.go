package resources

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecrpublic"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ECRPublicRepository struct {
	svc         *ecrpublic.ECRPublic
	name        *string
	createdTime *time.Time
	tags        []*ecrpublic.Tag
}

func init() {
	register("ECRPublicRepository", ListECRPublicRepositories)
}

func ListECRPublicRepositories(sess *session.Session) ([]Resource, error) {
	svc := ecrpublic.New(sess)
	resources := []Resource{}

	input := &ecrpublic.DescribeRepositoriesInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeRepositories(input)
		if err != nil {
			return nil, err
		}

		for _, repository := range output.Repositories {
			fmt.Println(repository)
			tagResp, err := svc.ListTagsForResource(&ecrpublic.ListTagsForResourceInput{
				ResourceArn: repository.RepositoryArn,
			})
			if err != nil {
				return nil, err
			}
			resources = append(resources, &ECRPublicRepository{
				svc:         svc,
				name:        repository.RepositoryName,
				createdTime: repository.CreatedAt,
				tags:        tagResp.Tags,
			})
		}

		if output.NextToken == nil {
			break
		}

		input.NextToken = output.NextToken
	}

	return resources, nil
}

func (r *ECRPublicRepository) Filter() error {
	return nil
}

func (r *ECRPublicRepository) Properties() types.Properties {
	properties := types.NewProperties().
		Set("CreatedTime", r.createdTime.Format(time.RFC3339))

	for _, t := range r.tags {
		properties.SetTag(t.Key, t.Value)
	}
	return properties
}

func (r *ECRPublicRepository) Remove() error {
	params := &ecrpublic.DeleteRepositoryInput{
		RepositoryName: r.name,
		Force:          aws.Bool(true),
	}
	_, err := r.svc.DeleteRepository(params)
	return err
}

func (r *ECRPublicRepository) String() string {
	return fmt.Sprintf("Repository: %s", *r.name)
}
