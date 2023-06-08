package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/gamelift"
)

type GameLiftMatchmakingRuleSet struct {
	svc  *gamelift.GameLift
	Name string
}

func init() {
	register("GameLiftMatchmakingRuleSet", ListMatchmakingRuleSets)
}

func ListMatchmakingRuleSets(sess *session.Session) ([]Resource, error) {
	svc := gamelift.New(sess)

	resp, err := svc.DescribeMatchmakingRuleSets(&gamelift.DescribeMatchmakingRuleSetsInput{})
	if err != nil {
		return nil, err
	}

	rules := make([]Resource, 0)
	for _, ruleSet := range resp.RuleSets {
		q := &GameLiftMatchmakingRuleSet{
			svc:  svc,
			Name: *ruleSet.RuleSetName,
		}
		rules = append(rules, q)
	}

	return rules, nil
}

func (ruleSet *GameLiftMatchmakingRuleSet) Remove() error {
	params := &gamelift.DeleteMatchmakingRuleSetInput{
		Name: aws.String(ruleSet.Name),
	}

	_, err := ruleSet.svc.DeleteMatchmakingRuleSet(params)
	if err != nil {
		return err
	}

	return nil
}

func (ruleSet *GameLiftMatchmakingRuleSet) String() string {
	return ruleSet.Name
}
