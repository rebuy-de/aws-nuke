package cmd

import "github.com/aws/aws-sdk-go/aws/session"

type Region struct {
	Name    string
	Session *session.Session
}
