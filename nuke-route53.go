package main

import "github.com/aws/aws-sdk-go/service/route53"

type Route53Nuke struct {
	svc *route53.Route53
}
