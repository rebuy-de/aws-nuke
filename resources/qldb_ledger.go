package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/qldb"
	"github.com/hunterkepley/aws-nuke/v2/pkg/config"
	"github.com/hunterkepley/aws-nuke/v2/pkg/types"
)

type QLDBLedger struct {
	svc    *qldb.QLDB
	ledger *qldb.DescribeLedgerOutput

	featureFlags config.FeatureFlags
}

func init() {
	register("QLDBLedger", ListQLDBLedgers)
}

func ListQLDBLedgers(sess *session.Session) ([]Resource, error) {
	svc := qldb.New(sess)

	params := &qldb.ListLedgersInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListLedgers(params)
		if err != nil {
			return nil, err
		}

		for _, ledger := range resp.Ledgers {
			ledgerDescription, err := svc.DescribeLedger(&qldb.DescribeLedgerInput{Name: ledger.Name})
			if err != nil {
				return nil, err
			}

			resources = append(resources, &QLDBLedger{
				svc:    svc,
				ledger: ledgerDescription,
			})
		}
		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (l *QLDBLedger) FeatureFlags(ff config.FeatureFlags) {
	l.featureFlags = ff
}

func (l *QLDBLedger) Remove() error {
	if aws.BoolValue(l.ledger.DeletionProtection) && l.featureFlags.DisableDeletionProtection.QLDBLedger {
		modifyParams := &qldb.UpdateLedgerInput{
			DeletionProtection: aws.Bool(false),
			Name:               l.ledger.Name,
		}
		_, err := l.svc.UpdateLedger((modifyParams))
		if err != nil {
			return err
		}
	}

	params := &qldb.DeleteLedgerInput{
		Name: l.ledger.Name,
	}

	_, err := l.svc.DeleteLedger(params)
	if err != nil {
		return err
	}

	return nil
}

func (l *QLDBLedger) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", l.ledger.Name)
	properties.Set("DeletionProtection", l.ledger.DeletionProtection)
	properties.Set("Arn", l.ledger.Arn)
	properties.Set("CreationDateTime", l.ledger.CreationDateTime.Format(time.RFC3339))
	properties.Set("State", l.ledger.State)
	properties.Set("PermissionsMode", l.ledger.PermissionsMode)
	properties.Set("EncryptionDescription", l.ledger.EncryptionDescription)
	return properties
}
func (l *QLDBLedger) String() string {
	return aws.StringValue(l.ledger.Name)
}
