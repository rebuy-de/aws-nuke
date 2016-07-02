package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/fatih/color"
)

var (
	// will be overwritten on build
	version = "unknown"
)

var (
	ReasonSkip            = *color.New(color.FgYellow)
	ReasonError           = *color.New(color.FgRed)
	ReasonRemoveTriggered = *color.New()
	ReasonWaitPending     = *color.New()
	ReasonSuccess         = *color.New(color.FgGreen)
	ColorID               = *color.New(color.Bold)
)

func Log(r Resource, c color.Color, msg string) {
	fmt.Printf("[%s] ", time.Now().Format(time.RFC3339))
	fmt.Print(strings.Split(fmt.Sprintf("%T", r), ".")[1]) // hackey
	fmt.Printf(" - ")
	ColorID.Printf("'%s'", r.String())
	fmt.Printf(" - ")
	c.Printf("%s\n", msg)
}

func main() {
	fmt.Printf("Running aws-nuke version %s.\n", version)

	var (
		noDryRun = flag.Bool("no-dry-run", false,
			"Actualy delete found resources.")
		noWait = flag.Bool("no-wait", false,
			"Do not wait for resource removal. This is faster, "+
				"but you may have to run the nuke multiple times.")
	)
	flag.Parse()

	if !*noDryRun {
		fmt.Printf("Dry run: do real delete with '--no-dry-run'.\n")
	}

	fmt.Println()

	nuke := &Nuke{
		session: session.New(&aws.Config{
			Region:      aws.String("eu-central-1"),
			Credentials: credentials.NewSharedCredentials("", "svenwltr"),
		}),
		dry:  !*noDryRun,
		wait: !*noWait,

		queue:    []Resource{},
		waiting:  []Resource{},
		skipped:  []Resource{},
		errored:  []Resource{},
		finished: []Resource{},
	}

	listers := nuke.GetListers()

	for _, lister := range listers {
		nuke.Scan(lister)
		nuke.CheckQueue()
		nuke.HandleQueue()
		nuke.Wait()
	}
}

type Nuke struct {
	session *session.Session

	dry  bool
	wait bool

	queue    []Resource
	waiting  []Resource
	skipped  []Resource
	errored  []Resource
	finished []Resource
}

func (n *Nuke) GetListers() []ResourceLister {
	ec2 := EC2Nuke{ec2.New(n.session)}
	autoscaling := AutoScalingNuke{autoscaling.New(n.session)}
	route53 := Route53Nuke{route53.New(n.session)}

	return []ResourceLister{
		autoscaling.ListGroups,
		ec2.ListInstances,
		ec2.ListSecurityGroups,
		ec2.ListVpcs,
		route53.ListResourceRecords,
		route53.ListHostedZones,
	}
}

func (n *Nuke) Scan(lister ResourceLister) error {
	resources, err := lister()
	if err != nil {
		return err
	}

	n.queue = append(n.queue, resources...)

	return nil
}

func (n *Nuke) CheckQueue() {
	temp := n.queue[:]
	n.queue = n.queue[0:0]

	for _, resource := range temp {
		checker, ok := resource.(Checker)
		if !ok {
			n.queue = append(n.queue, resource)
			continue
		}

		err := checker.Check()
		if err == nil {
			n.queue = append(n.queue, resource)
			continue
		}

		Log(resource, ReasonSkip, err.Error())
		n.skipped = append(n.skipped, resource)
	}
}

func (n *Nuke) HandleQueue() {
	temp := n.queue[:]
	n.queue = n.queue[0:0]

	for _, resource := range temp {
		if n.dry {
			n.skipped = append(n.skipped, resource)
			Log(resource, ReasonSuccess, "would remove")
			continue
		}

		err := resource.Remove()
		if err != nil {
			n.errored = append(n.errored, resource)
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
		n.waiting = []Resource{}
		return
	}

	temp := n.waiting[:]
	n.waiting = n.waiting[0:0]

	var wg sync.WaitGroup
	for i, resource := range temp {
		waiter, ok := resource.(Waiter)
		if !ok {
			continue
		}
		wg.Add(1)
		Log(resource, ReasonWaitPending, "waiting")
		go func(i int, resource Resource) {
			defer wg.Done()
			err := waiter.Wait()
			if err != nil {
				n.errored = append(n.errored, resource)
				Log(resource, ReasonError, err.Error())
				return
			}

			n.finished = append(n.finished, resource)
			Log(resource, ReasonSuccess, "removed")
		}(i, resource)
	}

	wg.Wait()

}
