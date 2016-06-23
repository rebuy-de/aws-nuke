package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/vendor/github.com/aws/aws-sdk-go/service/autoscaling"
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

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func nukeSession(sess *session.Session) {
	ec2svc := ec2.New(sess)
	autoscalingSvc := autoscaling.New(sess)

	nukeAutoScalingGroups(autoscalingSvc)
	nukeEC2Instances(ec2svc)
}
