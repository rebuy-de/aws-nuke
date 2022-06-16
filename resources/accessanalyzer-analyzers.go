package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/accessanalyzer"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AccessAnalyzer struct {
	svc    *accessanalyzer.AccessAnalyzer
	arn    string
	name   string
	status string
	tags   map[string]*string
}

func init() {
	register("AccessAnalyzer", ListAccessAnalyzer,
		mapCloudControl("AWS::AccessAnalyzer::Analyzer"))
}

func ListAccessAnalyzer(sess *session.Session) ([]Resource, error) {
	svc := accessanalyzer.New(sess)

	params := &accessanalyzer.ListAnalyzersInput{
		Type: aws.String("ACCOUNT"),
	}

	resources := make([]Resource, 0)
	err := svc.ListAnalyzersPages(params,
		func(page *accessanalyzer.ListAnalyzersOutput, lastPage bool) bool {
			for _, analyzer := range page.Analyzers {
				resources = append(resources, &AccessAnalyzer{
					svc:    svc,
					arn:    *analyzer.Arn,
					name:   *analyzer.Name,
					status: *analyzer.Status,
					tags:   analyzer.Tags,
				})
			}
			return true
		})
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (a *AccessAnalyzer) Remove() error {
	_, err := a.svc.DeleteAnalyzer(&accessanalyzer.DeleteAnalyzerInput{AnalyzerName: &a.name})

	return err
}

func (a *AccessAnalyzer) Properties() types.Properties {
	properties := types.NewProperties()

	properties.Set("ARN", a.arn)
	properties.Set("Name", a.name)
	properties.Set("Status", a.status)
	for k, v := range a.tags {
		properties.SetTag(&k, v)
	}

	return properties
}

func (a *AccessAnalyzer) String() string {
	return a.name
}
