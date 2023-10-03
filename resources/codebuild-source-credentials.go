package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CodeBuildSourceCredential struct {
	svc        *codebuild.CodeBuild
	Arn        *string
	AuthType   *string
	ServerType *string
}

func init() {
	register("CodeBuildSourceCredential", ListCodeBuildSourceCredential)
}

func ListCodeBuildSourceCredential(sess *session.Session) ([]Resource, error) {
	svc := codebuild.New(sess)
	resources := []Resource{}

	params := &codebuild.ListSourceCredentialsInput{}

	resp, err := svc.ListSourceCredentials(params)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}

	for _, credential := range resp.SourceCredentialsInfos {
		resources = append(resources, &CodeBuildSourceCredential{
			svc: svc,
			Arn: credential.Arn,
		})
	}

	return resources, nil
}

func (f *CodeBuildSourceCredential) Remove() error {
	_, err := f.svc.DeleteSourceCredentials(&codebuild.DeleteSourceCredentialsInput{Arn: f.Arn})
	return err
}

func (f *CodeBuildSourceCredential) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Arn", f.Arn)
	properties.Set("AuthType", f.AuthType)
	properties.Set("ServerType", f.ServerType)

	return properties
}
