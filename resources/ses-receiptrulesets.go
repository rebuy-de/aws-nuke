package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SESReceiptRuleSet struct {
	svc           *ses.SES
	name          *string
	activeRuleSet bool
}

func init() {
	register("SESReceiptRuleSet", ListSESReceiptRuleSets)
}

func ListSESReceiptRuleSets(sess *session.Session) ([]Resource, error) {
	svc := ses.New(sess)
	resources := []Resource{}

	params := &ses.ListReceiptRuleSetsInput{}

	output, err := svc.ListReceiptRuleSets(params)
	if err != nil {
		return nil, err
	}

	for _, ruleSet := range output.RuleSets {

		//Check active state
		ruleSetState := false
		ruleName := ruleSet.Name

		activeRuleSetOutput, err := svc.DescribeActiveReceiptRuleSet(&ses.DescribeActiveReceiptRuleSetInput{})
		if err != nil {
			return nil, err
		}
		if activeRuleSetOutput.Metadata == nil {
			ruleSetState = false
		} else if *ruleName == *activeRuleSetOutput.Metadata.Name {
			ruleSetState = true
		}

		resources = append(resources, &SESReceiptRuleSet{
			svc:           svc,
			name:          ruleName,
			activeRuleSet: ruleSetState,
		})
	}

	return resources, nil
}

func (f *SESReceiptRuleSet) Remove() error {

	_, err := f.svc.DeleteReceiptRuleSet(&ses.DeleteReceiptRuleSetInput{
		RuleSetName: f.name,
	})

	return err
}

func (f *SESReceiptRuleSet) String() string {
	return *f.name
}

func (f *SESReceiptRuleSet) Filter() error {
	if f.activeRuleSet == true {
		return fmt.Errorf("cannot delete active ruleset")
	}
	return nil
}
