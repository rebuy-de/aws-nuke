package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
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

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func nukeSession(sess *session.Session) {
	ec2svc := ec2.New(sess)
	nukeEC2Instances(ec2svc)
}

func nukeEC2Instances(ec2svc *ec2.EC2) {
	dii := &ec2.DescribeInstancesInput{}
	dio, err := ec2svc.DescribeInstances(dii)
	assertNoError(err)

	tii := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{},
	}

	log.Printf("Found these EC2 instances:")

	for _, reservation := range dio.Reservations {
		for _, instance := range reservation.Instances {
			log.Printf("\t%s (KeyName=%s, InstanceType=%s, State=%s)",
				*instance.InstanceId, *instance.KeyName,
				*instance.InstanceType, *instance.State.Name)

			if *instance.State.Name == "running" {
				tii.InstanceIds = append(tii.InstanceIds, aws.String(*instance.InstanceId))
			}
		}
	}

	if len(tii.InstanceIds) == 0 {
		log.Printf("Did not find any running instance.")
		return
	}

	log.Printf("Going to terminate these running EC2 instances: %v", tii.InstanceIds)
	_, err = ec2svc.TerminateInstances(tii)
	assertNoError(err)
}
