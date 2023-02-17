package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/gluedatabrew"
)

type GlueDataBrewRulesets struct {
	svc  *gluedatabrew.GlueDataBrew
	name *string
}

func init() {
	register("GlueDataBrewRulesets", ListGlueDataBrewRulesets)
}

func ListGlueDataBrewRulesets(sess *session.Session) ([]Resource, error) {
	svc := gluedatabrew.New(sess)
	resources := []Resource{}

	params := &gluedatabrew.ListRulesetsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListRulesets(params)
		if err != nil {
			return nil, err
		}

		for _, ruleset := range output.Rulesets {
			resources = append(resources, &GlueDataBrewRulesets{
				svc:  svc,
				name: ruleset.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueDataBrewRulesets) Remove() error {
	_, err := f.svc.DeleteRuleset(&gluedatabrew.DeleteRulesetInput{
		Name: f.name,
	})

	return err
}

func (f *GlueDataBrewRulesets) String() string {
	return *f.name
}
