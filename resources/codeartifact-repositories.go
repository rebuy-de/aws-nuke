package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codeartifact"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CodeArtifactRepository struct {
	svc    *codeartifact.CodeArtifact
	name   *string
	domain *string
	tags   map[string]*string
}

func init() {
	register("CodeArtifactRepository", ListCodeArtifactRepositories)
}

func ListCodeArtifactRepositories(sess *session.Session) ([]Resource, error) {
	svc := codeartifact.New(sess)
	resources := []Resource{}

	params := &codeartifact.ListRepositoriesInput{}

	for {
		resp, err := svc.ListRepositories(params)
		if err != nil {
			return nil, err
		}

		for _, repo := range resp.Repositories {
			resources = append(resources, &CodeArtifactRepository{
				svc:    svc,
				name:   repo.Name,
				domain: repo.DomainName,
				tags:   GetRepositoryTags(svc, repo.Arn),
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func GetRepositoryTags(svc *codeartifact.CodeArtifact, arn *string) map[string]*string {
	tags := map[string]*string{}

	resp, _ := svc.ListTagsForResource(&codeartifact.ListTagsForResourceInput{
		ResourceArn: arn,
	})
	for _, tag := range resp.Tags {
		tags[*tag.Key] = tag.Value
	}

	return tags
}

func (r *CodeArtifactRepository) Remove() error {
	_, err := r.svc.DeleteRepository(&codeartifact.DeleteRepositoryInput{
		Repository: r.name,
		Domain:     r.domain,
	})
	return err
}

func (r *CodeArtifactRepository) String() string {
	return *r.name
}

func (r *CodeArtifactRepository) Properties() types.Properties {
	properties := types.NewProperties()
	for key, tag := range r.tags {
		properties.SetTag(&key, tag)
	}
	properties.Set("Name", r.name)
	properties.Set("Domain", r.domain)
	return properties
}
