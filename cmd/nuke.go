package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/rebuy-de/aws-nuke/resources"
)

type Nuke struct {
	Parameters NukeParameters
	Config     *NukeConfig

	accountConfig NukeConfigAccount
	accountID     string
	accountAlias  string
	sessions      map[string]*session.Session

	ForceSleep time.Duration

	items Queue
}

func NewNuke(params NukeParameters) *Nuke {
	n := Nuke{
		Parameters: params,
		ForceSleep: 15 * time.Second,
	}

	return &n
}

func (n *Nuke) StartSession() error {
	n.sessions = make(map[string]*session.Session)
	for _, region := range n.Config.Regions {
		if n.Parameters.hasProfile() {
			s := session.New(&aws.Config{
				Region:      &region,
				Credentials: credentials.NewSharedCredentials("", n.Parameters.Profile),
			})

			if s == nil {
				return fmt.Errorf("Unable to create session with profile '%s'.", n.Parameters.Profile)
			}

			n.sessions[region] = s
			return nil
		}

		if n.Parameters.hasKeys() {
			s := session.New(&aws.Config{
				Region: &region,
				Credentials: credentials.NewStaticCredentials(
					n.Parameters.AccessKeyID,
					n.Parameters.SecretAccessKey,
					"",
				),
			})

			if s == nil {
				return fmt.Errorf("Unable to create session with key ID '%s'.", n.Parameters.AccessKeyID)
			}

			n.sessions[region] = s

		}
		return nil

	}

	return fmt.Errorf("You have to specify a profile or credentials.")
}

func (n *Nuke) Run() error {
	var err error

	fmt.Printf("aws-nuke version %s - %s - %s\n\n", BuildVersion, BuildDate, BuildHash)

	err = n.ValidateAccount()
	if err != nil {
		return err
	}

	fmt.Printf("Do you really want to nuke the account with "+
		"the ID %s and the alias '%s'?\n", n.accountID, n.accountAlias)
	if n.Parameters.Force {
		fmt.Printf("Waiting %v before continuing.\n", n.ForceSleep)
		time.Sleep(n.ForceSleep)
	} else {
		fmt.Printf("Do you want to continue? Enter account alias to continue.\n")
		err = Prompt(n.accountAlias)
		if err != nil {
			return err
		}
	}

	err = n.Scan()
	if err != nil {
		return err
	}

	if n.items.Count(ItemStateNew) == 0 {
		fmt.Println("No resource to delete.")
		return nil
	}

	if !n.Parameters.NoDryRun {
		fmt.Println("Would delete these resources. Provide --no-dry-run to actually destroy resources.")
		return nil
	}

	fmt.Printf("Do you really want to nuke these resources on the account with "+
		"the ID %s and the alias '%s'?\n", n.accountID, n.accountAlias)
	if n.Parameters.Force {
		fmt.Printf("Waiting %v before continuing.\n", n.ForceSleep)
		time.Sleep(n.ForceSleep)
	} else {
		fmt.Printf("Do you want to continue? Enter account alias to continue.\n")
		err = Prompt(n.accountAlias)
		if err != nil {
			return err
		}
	}

	failCount := 0

	for {
		n.HandleQueue()

		if n.items.Count(ItemStatePending, ItemStateWaiting, ItemStateNew) == 0 && n.items.Count(ItemStateFailed) > 0 {
			if failCount >= 2 {
				return fmt.Errorf("There are resources in failed state, but none are ready for deletion, anymore.")
			}
			failCount = failCount + 1
		} else {
			failCount = 0
		}

		if n.items.Count(ItemStateNew, ItemStatePending, ItemStateFailed, ItemStateWaiting) == 0 {
			break
		}

		time.Sleep(5 * time.Second)
	}

	fmt.Printf("Nuke complete: %d failed, %d skipped, %d finished.\n\n",
		n.items.Count(ItemStateFailed), n.items.Count(ItemStateFiltered), n.items.Count(ItemStateFinished))

	return nil
}

func (n *Nuke) ValidateAccount() error {
	sess := n.sessions[n.Config.Regions[0]]
	identOutput, err := sts.New(sess).GetCallerIdentity(nil)
	if err != nil {
		return err
	}

	aliasesOutput, err := iam.New(sess).ListAccountAliases(nil)
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
	queue := make(Queue, 0)

	for _, sess := range n.sessions {
		scanner := Scan(sess)
		for item := range scanner.Items {
			if !n.Parameters.WantsTarget(item.Service) {
				continue
			}

			queue = append(queue, item)
			n.Filter(item)
			item.Print()
		}

		if scanner.Error != nil {
			return scanner.Error
		}

	}
	fmt.Printf("Scan complete: %d total, %d nukeable, %d filtered.\n\n",
		queue.CountTotal(), queue.Count(ItemStateNew), queue.Count(ItemStateFiltered))

	n.items = queue

	return nil
}

func (n *Nuke) Filter(item *Item) {
	checker, ok := item.Resource.(resources.Filter)
	if ok {
		err := checker.Filter()
		if err != nil {
			item.State = ItemStateFiltered
			item.Reason = err.Error()
			return
		}
	}

	filters, ok := n.accountConfig.Filters[item.Service]
	if !ok {
		return
	}

	for _, filter := range filters {
		if filter == item.Resource.String() {
			item.State = ItemStateFiltered
			item.Reason = "filtered by config"
			return
		}
	}

	return
}

func (n *Nuke) HandleQueue() {
	listCache := make(map[string][]resources.Resource)

	for _, item := range n.items {
		switch item.State {
		case ItemStateNew:
			n.HandleRemove(item)
			item.Print()
		case ItemStateFailed:
			n.HandleRemove(item)
			n.HandleWait(item, listCache)
			item.Print()
		case ItemStatePending:
			n.HandleWait(item, listCache)
			item.State = ItemStateWaiting
			item.Print()
		case ItemStateWaiting:
			n.HandleWait(item, listCache)
			item.Print()
		}

	}

	fmt.Println()
	fmt.Printf("Removal requested: %d waiting, %d failed, %d skipped, %d finished\n\n",
		n.items.Count(ItemStateWaiting, ItemStatePending), n.items.Count(ItemStateFailed),
		n.items.Count(ItemStateFiltered), n.items.Count(ItemStateFinished))
}

func (n *Nuke) HandleRemove(item *Item) {
	err := item.Resource.Remove()
	if err != nil {
		item.State = ItemStateFailed
		item.Reason = err.Error()
		return
	}

	item.State = ItemStatePending
	item.Reason = ""
}

func (n *Nuke) HandleWait(item *Item, cache map[string][]resources.Resource) {
	var err error

	left, ok := cache[item.Service]
	if !ok {
		left, err = item.Lister()
		if err != nil {
			item.State = ItemStateFailed
			item.Reason = err.Error()
			return
		}
		cache[item.Service] = left
	}

	for _, r := range left {
		if r.String() == item.Resource.String() {
			checker, ok := r.(resources.Filter)
			if ok {
				err := checker.Filter()
				if err != nil {
					break
				}
			}

			return
		}
	}

	item.State = ItemStateFinished
	item.Reason = ""
}
