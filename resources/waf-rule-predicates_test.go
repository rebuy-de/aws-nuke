package resources

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/golang/mock/gomock"
	"github.com/rebuy-de/aws-nuke/mocks/mock_wafiface"
	"github.com/stretchr/testify/assert"
)

func TestWAFPredicates_Remove(t *testing.T) {
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWAFRulePredicate := mock_wafiface.NewMockWAFAPI(ctrl)

	wafRulePredicate := WAFRulePredicate{
		svc:    mockWAFRulePredicate,
		ruleID: aws.String("foobar"),
		predicate: &waf.Predicate{
			Type:    aws.String("IPSet"),
			Negated: aws.Bool(false),
			DataId:  aws.String("foobar"),
		},
	}

	gomock.InOrder(
		mockWAFRulePredicate.EXPECT().GetChangeToken(gomock.Eq(&waf.GetChangeTokenInput{})).
			Return(&waf.GetChangeTokenOutput{
				ChangeToken: aws.String("foobar")}, nil),
		mockWAFRulePredicate.EXPECT().UpdateRule(gomock.Eq(&waf.UpdateRuleInput{
			ChangeToken: aws.String("foobar"),
			RuleId:      aws.String("foobar"),
			Updates: []*waf.RuleUpdate{
				&waf.RuleUpdate{
					Action:    aws.String("DELETE"),
					Predicate: wafRulePredicate.predicate,
				},
			},
		})).Return(nil, nil),
	)

	err := wafRulePredicate.Remove()
	a.Nil(err)
}
