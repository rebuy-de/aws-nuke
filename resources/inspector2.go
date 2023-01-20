package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/inspector2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type Inspector2 struct {
	svc       *inspector2.Inspector2
	accountId *string
}

func init() {
	register("Inspector2", ListInspector2)
}

func ListInspector2(sess *session.Session) ([]Resource, error) {
	svc := inspector2.New(sess)

	resources := make([]Resource, 0)

	resp, err := svc.BatchGetAccountStatus(nil)
	if err != nil {
		return resources, err
	}
	for _, a := range resp.Accounts {
		if *a.State.Status != inspector2.StatusDisabled {
			resources = append(resources, &Inspector2{
				svc:       svc,
				accountId: a.AccountId,
			})
		}
	}

	return resources, nil
}

func (e *Inspector2) Remove() error {
	_, err := e.svc.Disable(&inspector2.DisableInput{
		AccountIds:    []*string{e.accountId},
		ResourceTypes: aws.StringSlice(inspector2.ResourceScanType_Values()),
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *Inspector2) Properties() types.Properties {
	properties := types.NewProperties()

	properties.Set("AccountID", e.accountId)

	return properties
}

func (e *Inspector2) String() string {
	return *e.accountId
}
