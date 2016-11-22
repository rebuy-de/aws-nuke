package cmd

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/rebuy-de/aws-nuke/resources"
)

type Nuke struct {
	Parameters NukeParameters

	Config  *NukeConfig
	session *session.Session

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

	err = n.ValidateAccount()
	if err != nil {
		return err
	}

	err = n.Scan()
	if err != nil {
		return err
	}

	fmt.Printf("\nScan complete: %d total, %d nukeable, %d filtered.\n\n",
		len(n.queue)+len(n.skipped), len(n.queue), len(n.skipped))

	return nil

	n.HandleQueue()
	n.Wait()

	if n.retry {
		for len(n.failed) > 0 {
			fmt.Println()
			fmt.Printf("Retrying: %d finished, %d failed, %d skipped.",
				len(n.finished), len(n.failed), len(n.skipped))
			fmt.Println()
			fmt.Println()
			n.Retry()
		}
	}

	fmt.Println()
	fmt.Printf("Nuke complete: %d finished, %d failed, %d skipped.",
		len(n.finished), len(n.failed), len(n.skipped))
	fmt.Println()

	return err
}

func (n *Nuke) ValidateAccount() error {
	ident, err := sts.New(n.session).GetCallerIdentity(nil)
	if err != nil {
		return err
	}
	accountID := *ident.Account

	if !n.Config.HasBlacklist() {
		return fmt.Errorf("The config file contains an empty blacklist. " +
			"For safety reasons you need to specify at least one account ID. " +
			"This should be you production account.")
	}

	if n.Config.InBlacklist(accountID) {
		return fmt.Errorf("You are trying to nuke the account with the ID %s, "+
			"but it is blacklisted. Aborting.", accountID)
	}

	return AskContinue("Do you really want to nuke the account with the ID %s?", accountID)
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

	return nil
}

func (n *Nuke) Retry() {
	n.queue = n.failed[:]
	n.failed = n.failed[0:0]

	n.HandleQueue()
	n.Wait()
}

func (n *Nuke) HandleQueue() {
	temp := n.queue[:]
	n.queue = n.queue[0:0]

	for _, resource := range temp {
		if !n.Parameters.NoDryRun {
			n.skipped = append(n.skipped, resource)
			continue
		}

		err := resource.Remove()
		if err != nil {
			n.failed = append(n.failed, resource)
			Log(resource, ReasonError, err.Error())
			continue
		}

		n.waiting = append(n.waiting, resource)
		Log(resource, ReasonRemoveTriggered, "triggered remove")
	}
}

func (n *Nuke) Wait() {
	if !n.wait {
		n.finished = n.waiting
		n.waiting = []resources.Resource{}
		return
	}

	temp := n.waiting[:]
	n.waiting = n.waiting[0:0]

	var wg sync.WaitGroup
	for i, resource := range temp {
		waiter, ok := resource.(resources.Waiter)
		if !ok {
			n.finished = append(n.finished, resource)
			continue
		}
		wg.Add(1)
		Log(resource, ReasonWaitPending, "waiting")
		go func(i int, resource resources.Resource) {
			defer wg.Done()
			err := waiter.Wait()
			if err != nil {
				n.failed = append(n.failed, resource)
				Log(resource, ReasonError, err.Error())
				return
			}

			n.finished = append(n.finished, resource)
			Log(resource, ReasonSuccess, "removed")
		}(i, resource)
	}

	wg.Wait()
}
