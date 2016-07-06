package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Nuke struct {
	session *session.Session

	dry       bool
	wait      bool
	earlyExit bool
	retry     bool

	queue    []Resource
	waiting  []Resource
	skipped  []Resource
	failed   []Resource
	finished []Resource
}

func (n *Nuke) Run() {
	listers := n.GetListers()

	for _, lister := range listers {
		err := n.Scan(lister)
		if err != nil {
			LogErrorf(err)
			continue
		}

		n.CheckQueue()
		n.HandleQueue()
		n.Wait()
	}

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
}

func (n *Nuke) GetListers() []ResourceLister {
	autoscaling := AutoScalingNuke{autoscaling.New(n.session)}
	ec2 := EC2Nuke{ec2.New(n.session)}
	elb := ElbNuke{elb.New(n.session)}
	route53 := Route53Nuke{route53.New(n.session)}
	s3 := S3Nuke{s3.New(n.session)}
	iam := IamNuke{iam.New(n.session)}

	return []ResourceLister{
		elb.ListELBs,
		autoscaling.ListGroups,
		route53.ListResourceRecords,
		route53.ListHostedZones,
		ec2.ListKeyPairs,
		ec2.ListInstances,
		ec2.ListSecurityGroups,
		ec2.ListVpnGatewayAttachements,
		ec2.ListVpnConnections,
		ec2.ListNetworkACLs,
		ec2.ListSubnets,
		ec2.ListCustomerGateways,
		ec2.ListVpnGateways,
		ec2.ListInternetGatewayAttachements,
		ec2.ListInternetGateways,

		ec2.ListRouteTables,
		ec2.ListDhcpOptions,
		ec2.ListVpcs,

		iam.ListInstanceProfileRoles,
		iam.ListInstanceProfiles,
		iam.ListRolePolicyAttachements,
		iam.ListRoles,

		s3.ListObjects,
		s3.ListBuckets,
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
		if n.dry {
			n.skipped = append(n.skipped, resource)
			Log(resource, ReasonSuccess, "would remove")
			continue
		}

		err := resource.Remove()
		if err != nil {
			n.failed = append(n.failed, resource)
			Log(resource, ReasonError, err.Error())
			if n.earlyExit {
				os.Exit(1)
			} else {
				continue
			}
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
			n.finished = append(n.finished, resource)
			continue
		}
		wg.Add(1)
		Log(resource, ReasonWaitPending, "waiting")
		go func(i int, resource Resource) {
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
