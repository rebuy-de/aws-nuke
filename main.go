package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	// will be overwritten on build
	version = "unknown"
)

func main() {
	fmt.Printf("Running aws-nuke version %s.\n", version)

	var (
		noDryRun = flag.Bool("no-dry-run", false,
			"actualy delete found resources")
		noWait = flag.Bool("no-wait", false,
			"do not wait for resource removal")
	)
	flag.Parse()

	if !*noDryRun {
		fmt.Printf("Dry running nuke. Do real delete with '--no-dry-run'.\n")
	}

	fmt.Println()

	credentials := credentials.NewSharedCredentials("", "svenwltr")
	sess := session.New(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: credentials,
	})

	nukeSession(sess, !*noDryRun, !*noWait)
}

func nukeSession(sess *session.Session, dry bool, wait bool) {
	ec2Nuke := EC2Nuke{ec2.New(sess)}
	autoscalingNuke := AutoScalingNuke{autoscaling.New(sess)}

	listers := []ResourceLister{
		autoscalingNuke.ListGroups,
		ec2Nuke.ListInstances,
		ec2Nuke.ListSecurityGroups,
	}

	for _, lister := range listers {
		err := nukeResource(lister, dry, wait)

		if err != nil {
			fmt.Fprintf(os.Stderr, "\n%s\n", err.Error())
		}
	}
}

func nukeResource(lister ResourceLister, dry bool, wait bool) error {
	var resources []Resource
	var err error

	resources, err = lister()
	if err != nil {
		return err
	}

	queue := make([]Resource, 0)
	for _, resource := range resources {
		fmt.Printf("%T %s", resource, resource.String())

		err = resource.Check()
		if err != nil {
			fmt.Printf(" ... %s\n", err.Error())
			continue
		}

		if dry {
			fmt.Printf(" ... would be removed\n")
			continue
		}

		err = resource.Remove()
		if err != nil {
			fmt.Printf(" ... %s\n", err.Error())
			continue
		}

		fmt.Printf(" ... delete requested [%d]\n", len(queue))
		queue = append(queue, resource)
	}

	if wait && len(queue) > 0 {
		fmt.Printf("Waiting, until %d resources get removed.", len(queue))
		var wg sync.WaitGroup
		for i, resource := range queue {
			wg.Add(1)
			go func(i int, resource Resource) {
				defer wg.Done()
				resource.Wait()
				fmt.Printf("%T %s ... deleted\n", resource, resource.String())
			}(i, resource)
		}

		wg.Wait()
	}

	return nil
}
