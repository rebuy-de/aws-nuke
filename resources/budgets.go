package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/budgets"
	"github.com/aws/aws-sdk-go/service/sts"
)

func init() {
	register("Budgets", ListBudgets)
}

type Budget struct {
	svc  *budgets.Budgets
	name *string
	// tags []*budgets.Tag
}

// if there are no tags what should properties be?

func ListBudgets(sess *session.Session) ([]Resource, error) {
	svc := budgets.New(sess)

	resources := []Resource{}

	// Lookup current account ID
	fmt.Println("Entered Budgets Resource")

	stsSvc := sts.New(sess)

	callerID, err := stsSvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	accountID := callerID.Account
	// fmt.Printf("account_num: %s \n", *accountID)

	params := &budgets.DescribeBudgetsInput{
		AccountId: aws.String(*accountID),
	}

	output, err := svc.DescribeBudgets(params)
	if err != nil {
		return nil, err
	}

	for _, bud := range output.Budgets {
		fmt.Println(bud.BudgetName)
		resources = append(resources, &Budget{
			svc:  svc,
			name: bud.BudgetName,
		})
	}

	return nil, nil
}
