package awsutil

import (
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
)

type Account struct {
	Credentials

	id      string
	aliases []string
}

func NewAccount(creds Credentials) (*Account, error) {
	account := Account{
		Credentials: creds,
	}

	defaultSession, err := account.NewSession(DefaultRegionID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create default session in %s", DefaultRegionID)
	}

	identityOutput, err := sts.New(defaultSession).GetCallerIdentity(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed get caller identity")
	}

	globalSession, err := account.NewSession(GlobalRegionID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create global session in %s", GlobalRegionID)
	}

	aliasesOutput, err := iam.New(globalSession).ListAccountAliases(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed get account alias")
	}

	aliases := []string{}
	for _, alias := range aliasesOutput.AccountAliases {
		if alias != nil {
			aliases = append(aliases, *alias)
		}
	}

	account.id = *identityOutput.Account
	account.aliases = aliases

	return &account, nil
}

func (a *Account) ID() string {
	return a.id
}

func (a *Account) Alias() string {
	return a.aliases[0]
}

func (a *Account) Aliases() []string {
	return a.aliases
}
