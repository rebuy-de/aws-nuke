package cmd

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/rebuy-de/aws-nuke/cmd/version"
	"github.com/rebuy-de/aws-nuke/resources"
)

type Nuke struct {
	Parameters NukeParameters
	Config     *NukeConfig

	accountConfig NukeConfigAccount
	accountID     string
	accountAlias  string
	session       *session.Session

	retry bool
	wait  bool

	queue    []resources.Resource
	waiting  []resources.Resource
	skipped  []resources.Resource
	failed   []resources.Resource
	finished []resources.Resource
}

func NewNuke(params NukeParameters) *Nuke {
	n := Nuke{
		Parameters: params,

		retry: true,
		wait:  true,

		queue:    []resources.Resource{},
		waiting:  []resources.Resource{},
		skipped:  []resources.Resource{},
		failed:   []resources.Resource{},
		finished: []resources.Resource{},
	}

	return &n
}

func (n *Nuke) StartSession() error {
	if n.Parameters.hasProfile() {
		s := session.New(&aws.Config{
			Region:      &n.Config.Region,
			Credentials: credentials.NewSharedCredentials("", n.Parameters.Profile),
		})

		if s == nil {
			return fmt.Errorf("Unable to create session with profile '%s'.", n.Parameters.Profile)
		}

		n.session = s
		return nil
	}

	if n.Parameters.hasKeys() {
		s := session.New(&aws.Config{
			Region: &n.Config.Region,
			Credentials: credentials.NewStaticCredentials(
				n.Parameters.AccessKeyID,
				n.Parameters.SecretAccessKey,
				"",
			),
		})

		if s == nil {
			return fmt.Errorf("Unable to create session with key ID '%s'.", n.Parameters.AccessKeyID)
		}

		n.session = s
		return nil
	}

	return fmt.Errorf("You have to specify a profile or credentials.")
}

func (n *Nuke) Run() error {
	var err error

	version.Print()

	err = n.ValidateAccount()
	if err != nil {
		return err
	}

	err = AskContinue("Do you really want to nuke the account with "+
		"the ID %s and the alias '%s'?", n.accountID, n.accountAlias)
	if err != nil {
		return err
	}

	err = n.Scan()
	if err != nil {
		return err
	}

	fmt.Printf("\nScan complete: %d total, %d nukeable, %d filtered.\n\n",
		len(n.queue)+len(n.skipped), len(n.queue), len(n.skipped))

	if len(n.queue) == 0 {
		fmt.Println("No resource to delete.")
		return nil
	}

	err = AskContinue("Do you really want to nuke these resources on the account with "+
		"the ID %s and the alias '%s'?", n.accountID, n.accountAlias)
	if err != nil {
		return err
	}

	for len(n.queue) != 0 {

		n.NukeQueue()
		n.WaitQueue()

		fmt.Println()
		fmt.Printf("Removal requested: %d failed, %d skipped, %d finished",
			len(n.failed), len(n.skipped), len(n.finished))
		fmt.Println()

		n.queue = n.failed
		n.failed = []resources.Resource{}

		time.Sleep(5 * time.Second)
	}

	fmt.Println()
	fmt.Printf("Nuke complete: %d failed, %d skipped, %d finished.",
		len(n.failed), len(n.skipped), len(n.finished))
	fmt.Println()

	return nil
}

func (n *Nuke) ValidateAccount() error {
	identOutput, err := sts.New(n.session).GetCallerIdentity(nil)
	if err != nil {
		return err
	}

	aliasesOutput, err := iam.New(n.session).ListAccountAliases(nil)
	if err != nil {
		return err
	}

	accountID := *identOutput.Account
	aliases := aliasesOutput.AccountAliases

	if !n.Config.HasBlacklist() {
		return fmt.Errorf("The config file contains an empty blacklist. " +
			"For safety reasons you need to specify at least one account ID. " +
			"This should be you production account.")
	}

	if n.Config.InBlacklist(accountID) {
		return fmt.Errorf("You are trying to nuke the account with the ID %s, "+
			"but it is blacklisted. Aborting.", accountID)
	}

	if len(aliases) == 0 {
		return fmt.Errorf("The specified account doesn't have an alias. " +
			"For safety reasons you need to specify an account alias. " +
			"Your production account should contain the term 'prod'.")
	}

	for _, alias := range aliases {
		if strings.Contains(strings.ToLower(*alias), "prod") {
			return fmt.Errorf("You are trying to nuke a account with the alias '%s', "+
				"but it has the substring 'prod' in it. Aborting.", *aliases[0])
		}
	}

	if _, ok := n.Config.Accounts[accountID]; !ok {
		return fmt.Errorf("Your account ID '%s' isn't listed in the config. "+
			"Aborting.", accountID)
	}

	n.accountConfig = n.Config.Accounts[accountID]
	n.accountID = accountID
	n.accountAlias = *aliases[0]

	return nil
}

func (n *Nuke) Scan() error {
	listers := resources.GetListers(n.session)

	for _, lister := range listers {
		resources, err := lister()
		if err != nil {
			return err
		}

		for _, r := range resources {
			reason := n.CheckFilters(r)
			if reason != nil {
				Log(r, ReasonSkip, reason.Error())
				n.skipped = append(n.skipped, r)
				continue
			}

			Log(r, ReasonSuccess, "would remove")
			n.queue = append(n.queue, r)
		}

	}

	return nil
}

func (n *Nuke) CheckFilters(r resources.Resource) error {
	checker, ok := r.(resources.Filter)
	if ok {
		err := checker.Filter()
		if err != nil {
			return err
		}
	}

	cat := resources.GetCategory(r)
	name := r.String()

	filters, ok := n.accountConfig.Filters[cat]
	if !ok {
		return nil
	}

	for _, filter := range filters {
		if filter == name {
			return fmt.Errorf("filtered by config")
		}
	}

	return nil
}

func (n *Nuke) NukeQueue() {
	for _, resource := range n.queue {
		err := resource.Remove()
		if err != nil {
			n.failed = append(n.failed, resource)
			Log(resource, ReasonError, err.Error())
			continue
		}

		n.waiting = append(n.waiting, resource)
		Log(resource, ReasonRemoveTriggered, "triggered remove")
	}

	n.queue = []resources.Resource{}
}

func (n *Nuke) WaitQueue() {
	var wg sync.WaitGroup

	for _, resource := range n.waiting {
		waiter, ok := resource.(resources.Waiter)
		if !ok {
			n.finished = append(n.finished, resource)
			Log(resource, ReasonSuccess, "deleted")
			continue
		}

		wg.Add(1)
		Log(resource, ReasonWaitPending, "waiting")

		go func(resource resources.Resource) {
			defer wg.Done()
			err := waiter.Wait()
			if err != nil {
				n.failed = append(n.failed, resource)
				Log(resource, ReasonError, err.Error())
				return
			}

			n.finished = append(n.finished, resource)
			Log(resource, ReasonSuccess, "removed")
		}(resource)
	}

	n.waiting = []resources.Resource{}
}
