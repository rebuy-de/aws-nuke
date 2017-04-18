package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type ECRrepository struct {
	svc    *ecr.ECR
	name   *string
	region *string
}

func (n *ECRNuke) ListRepos() ([]Resource, error) {

	var params *ecr.DescribeRepositoriesInput
	var resp *ecr.DescribeRepositoriesOutput
	var resources []Resource
	var err error
	for moreRepositories := true; moreRepositories; {
		if resp == nil {
			params = &ecr.DescribeRepositoriesInput{
				MaxResults: aws.Int64(100),
			}
			moreRepositories = true
		} else {
			if resp.NextToken != nil {
				fmt.Printf("Next token\n")
				params = &ecr.DescribeRepositoriesInput{
					MaxResults: aws.Int64(100),
					NextToken:  resp.NextToken,
				}
				moreRepositories = true
			} else {
				moreRepositories = false
				continue
			}
		}
		resp, err = n.Service.DescribeRepositories(params)
		if err != nil {
			fmt.Printf("Error occured")
			return nil, err
		}
		for _, repository := range resp.Repositories {
			resources = append(resources, &ECRrepository{
				svc:    n.Service,
				name:   repository.RepositoryName,
				region: n.Service.Config.Region,
			})
		}
	}

	return resources, nil
}

func (r *ECRrepository) Filter() error {
	return nil
}

func (r *ECRrepository) Remove() error {

	params := &ecr.DeleteRepositoryInput{
		RepositoryName: r.name,
		Force:          aws.Bool(true),
	}
	_, err := r.svc.DeleteRepository(params)
	return err
}

func (r *ECRrepository) String() string {
	return fmt.Sprintf("Repository: %s in %s", *r.name, *r.region)
}
