package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codecommit"
)

type CodeCommitRepository struct {
	svc            *codecommit.CodeCommit
	repositoryName *string
}

func init() {
	register("CodeCommitRepository", ListCodeCommitRepositories)
}

func ListCodeCommitRepositories(sess *session.Session) ([]Resource, error) {
	svc := codecommit.New(sess)
	resources := []Resource{}

	params := &codecommit.ListRepositoriesInput{}

	for {
		resp, err := svc.ListRepositories(params)
		if err != nil {
			return nil, err
		}

		for _, repository := range resp.Repositories {
			resources = append(resources, &CodeCommitRepository{
				svc:            svc,
				repositoryName: repository.RepositoryName,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *CodeCommitRepository) Remove() error {

	_, err := f.svc.DeleteRepository(&codecommit.DeleteRepositoryInput{
		RepositoryName: f.repositoryName,
	})

	return err
}

func (f *CodeCommitRepository) String() string {
	return *f.repositoryName
}
