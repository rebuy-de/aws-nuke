package resources

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ECRRepository struct {
	svc         *ecr.ECR
	name        *string
	createdTime *time.Time
	tags        []*ecr.Tag
}

func init() {
	register("ECRRepository", ListECRRepositories,
		mapCloudControl("AWS::ECR::Repository"))
}

func ListECRRepositories(sess *session.Session) ([]Resource, error) {
	svc := ecr.New(sess)
	resources := []Resource{}

	input := &ecr.DescribeRepositoriesInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeRepositories(input)
		if err != nil {
			return nil, err
		}

		for _, repository := range output.Repositories {
			tagResp, err := svc.ListTagsForResource(&ecr.ListTagsForResourceInput{
				ResourceArn: repository.RepositoryArn,
			})
			if err != nil {
				return nil, err
			}
			resources = append(resources, &ECRRepository{
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

func (r *ECRRepository) Filter() error {
	return nil
}

func (r *ECRRepository) Properties() types.Properties {
	properties := types.NewProperties().
		Set("CreatedTime", r.createdTime.Format(time.RFC3339))

	for _, t := range r.tags {
		properties.SetTag(t.Key, t.Value)
	}
	return properties
}

func (r *ECRRepository) Remove() error {
	params := &ecr.DeleteRepositoryInput{
		RepositoryName: r.name,
		Force:          aws.Bool(true),
	}
	_, err := r.svc.DeleteRepository(params)
	return err
}

func (r *ECRRepository) String() string {
	return fmt.Sprintf("Repository: %s", *r.name)
}
