package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/accessanalyzer"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ArchiveRule struct {
	svc          *accessanalyzer.AccessAnalyzer
	ruleName     string
	analyzerName string
}

func init() {
	register("ArchiveRule", ListArchiveRule)
}

func ListArchiveRule(sess *session.Session) ([]Resource, error) {
	svc := accessanalyzer.New(sess)

	analyzers, err := ListAccessAnalyzer(sess)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)

	for _, analyzer := range analyzers {
		a, ok := analyzer.(*AccessAnalyzer)
		if !ok {
			continue
		}

		params := &accessanalyzer.ListArchiveRulesInput{
			AnalyzerName: &a.name,
		}

		err = svc.ListArchiveRulesPages(params,
			func(page *accessanalyzer.ListArchiveRulesOutput, lastPage bool) bool {
				for _, archiveRule := range page.ArchiveRules {
					resources = append(resources, &ArchiveRule{
						svc:          svc,
						ruleName:     *archiveRule.RuleName,
						analyzerName: a.name,
					})
				}
				return true
			})
		if err != nil {
			return nil, err
		}
	}

	return resources, nil
}

func (a *ArchiveRule) Remove() error {
	_, err := a.svc.DeleteArchiveRule(&accessanalyzer.DeleteArchiveRuleInput{
		AnalyzerName: &a.analyzerName,
		RuleName:     &a.ruleName,
	})

	return err
}

func (a *ArchiveRule) Properties() types.Properties {
	properties := types.NewProperties()

	properties.Set("RuleName", a.ruleName)
	properties.Set("AnalyzerName", a.analyzerName)

	return properties
}

func (a *ArchiveRule) String() string {
	return a.ruleName
}
