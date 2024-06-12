package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CodeBuildReportGroup struct {
	svc *codebuild.CodeBuild
	Arn *string
}

func init() {
	register("CodeBuildReportGroup", ListCodeBuildReportGroup)
}

func ListCodeBuildReportGroup(sess *session.Session) ([]Resource, error) {
	svc := codebuild.New(sess)
	resources := []Resource{}

	params := &codebuild.ListReportGroupsInput{}

	for {
		resp, err := svc.ListReportGroups(params)
		if err != nil {
			return nil, err
		}

		for _, reportGroup := range resp.ReportGroups {
			resources = append(resources, &CodeBuildReportGroup{
				svc: svc,
				Arn: reportGroup,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *CodeBuildReportGroup) Remove() error {
	_, err := f.svc.DeleteReportGroup(&codebuild.DeleteReportGroupInput{
		Arn: f.Arn,
	})

	return err
}

func (f *CodeBuildReportGroup) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("Arn", f.Arn)
	return properties
}
