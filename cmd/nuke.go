package cmd

import (
	"fmt"
	"time"

	"github.com/rebuy-de/aws-nuke/pkg/awsutil"
	"github.com/rebuy-de/aws-nuke/resources"
)

type Nuke struct {
	Parameters NukeParameters
	Account    awsutil.Account
	Config     *NukeConfig

	accountConfig NukeConfigAccount

	ForceSleep time.Duration

	items Queue
}

func NewNuke(params NukeParameters, account awsutil.Account) *Nuke {
	n := Nuke{
		Parameters: params,
		Account:    account,
		ForceSleep: 15 * time.Second,
	}

	return &n
}

func (n *Nuke) Run() error {
	var err error

	fmt.Printf("aws-nuke version %s - %s - %s\n\n", BuildVersion, BuildDate, BuildHash)

	accountConfig, err := n.Config.ValidateAccount(n.Account.ID(), n.Account.Aliases())
	if err != nil {
		return err
	}
	n.accountConfig = *accountConfig

	fmt.Printf("Do you really want to nuke the account with "+
		"the ID %s and the alias '%s'?\n", n.Account.ID(), n.Account.Alias())
	if n.Parameters.Force {
		fmt.Printf("Waiting %v before continuing.\n", n.ForceSleep)
		time.Sleep(n.ForceSleep)
	} else {
		fmt.Printf("Do you want to continue? Enter account alias to continue.\n")
		err = Prompt(n.Account.Alias())
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
		"the ID %s and the alias '%s'?\n", n.Account.ID(), n.Account.Alias())
	if n.Parameters.Force {
		fmt.Printf("Waiting %v before continuing.\n", n.ForceSleep)
		time.Sleep(n.ForceSleep)
	} else {
		fmt.Printf("Do you want to continue? Enter account alias to continue.\n")
		err = Prompt(n.Account.Alias())
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

func (n *Nuke) Scan() error {
	queue := make(Queue, 0)

	for _, region := range n.Config.Regions {
		sess, err := n.Account.Session(region)
		if err != nil {
			return err
		}

		items := Scan(sess)
		for item := range items {
			if !n.Parameters.WantsTarget(item.Service) {
				continue
			}

			queue = append(queue, item)
			n.Filter(item)
			item.Print()
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
