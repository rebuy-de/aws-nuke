package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type ECRRepository struct {
	svc  *ecr.ECR
	name *string
}

func init() {
	register("ECRRepository", ListECRRepositories)
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
			resources = append(resources, &ECRRepository{
				svc:  svc,
				name: repository.RepositoryName,
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
