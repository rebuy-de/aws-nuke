package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/budgets"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func init() {
	register("Budget", ListBudgets)
}

type Budget struct {
	svc        *budgets.Budgets
	name       *string
	budgetType *string
	accountId  *string
}

func ListBudgets(sess *session.Session) ([]Resource, error) {
	svc := budgets.New(sess)

	identityOutput, err := sts.New(sess).GetCallerIdentity(nil)
	if err != nil {
		fmt.Printf("sts error: %s \n", err)
		return nil, err
	}
	accountID := identityOutput.Account

	params := &budgets.DescribeBudgetsInput{
		AccountId:  aws.String(*accountID),
		MaxResults: aws.Int64(100),
	}

	buds := make([]*budgets.Budget, 0)
	err = svc.DescribeBudgetsPages(params, func(page *budgets.DescribeBudgetsOutput, lastPage bool) bool {
		for _, out := range page.Budgets {
			buds = append(buds, out)
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	resources := []Resource{}
	for _, bud := range buds {
		resources = append(resources, &Budget{
			svc:        svc,
			name:       bud.BudgetName,
			budgetType: bud.BudgetType,
			accountId:  accountID,
		})
	}

	return resources, nil
}

func (b *Budget) Remove() error {

	_, err := b.svc.DeleteBudget(&budgets.DeleteBudgetInput{
		AccountId:  b.accountId,
		BudgetName: b.name,
	})

	return err
}

func (b *Budget) Properties() types.Properties {
	properties := types.NewProperties()

	properties.
		Set("Name", *b.name).
		Set("BudgetType", *b.budgetType).
		Set("AccountID", *b.accountId)
	return properties
}

func (b *Budget) String() string {
	return *b.name
}
