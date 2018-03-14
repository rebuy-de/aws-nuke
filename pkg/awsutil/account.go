package awsutil

import (
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
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

	defaultSession, err := account.Session(DefaultRegionID)
	if err != nil {
		return nil, err
	}

	identityOutput, err := sts.New(defaultSession).GetCallerIdentity(nil)
	if err != nil {
		return nil, err
	}

	globalSession, err := account.Session(GlobalRegionID)
	if err != nil {
		return nil, err
	}

	aliasesOutput, err := iam.New(globalSession).ListAccountAliases(nil)
	if err != nil {
		return nil, err
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
