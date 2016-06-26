package main

import (
	"fmt"
	"log"
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
	log.Printf("Running aws-nuke version %s.", version)

	credentials := credentials.NewSharedCredentials("", "svenwltr")
	sess := session.New(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: credentials,
	})

	nukeSession(sess)
}

func nukeSession(sess *session.Session) {
	ec2Nuke := EC2Nuke{ec2.New(sess)}
	autoscalingNuke := AutoScalingNuke{autoscaling.New(sess)}

	nukeResource(autoscalingNuke.ListGroups)
	nukeResource(ec2Nuke.ListInstances)
}

func nukeResource(lister ResourceLister) error {
	var resources []Resource
	var err error

	resources, err = lister()
	if err != nil {
		return err
	}

	for i, resource := range resources {
		fmt.Printf("%T %s", resource, resource.String())
		err = resource.Remove()
		if err != nil {
			return err
		}
		fmt.Printf(" [%d]\n", i)
	}

	if len(resources) > 0 {
		fmt.Printf("Waiting for %d resources: ", len(resources))
		var wg sync.WaitGroup

		for i, resource := range resources {
			wg.Add(1)
			go func(i int, resource Resource) {
				defer wg.Done()
				resource.Wait()
				fmt.Printf("[%d] ", i)
			}(i, resource)
		}

		wg.Wait()
		fmt.Println()
	}

	fmt.Println()

	return nil
}
