package awsutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
)

type Credentials struct {
	Profile string

	AccessKeyID     string
	SecretAccessKey string

	cache map[string]*session.Session
}

func (c *Credentials) HasProfile() bool {
	return strings.TrimSpace(c.Profile) != ""
}

func (c *Credentials) HasKeys() bool {
	return strings.TrimSpace(c.AccessKeyID) != "" &&
		strings.TrimSpace(c.SecretAccessKey) != ""
}

func (c *Credentials) Validate() error {
	if c.HasProfile() == c.HasKeys() {
		return fmt.Errorf("You have to specify the --profile flag OR " +
			"--access-key-id and --secret-access-key.\n")
	}

	return nil
}

func (c *Credentials) NewSession(region string) (*session.Session, error) {
	if c.HasProfile() {
		return Log(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region: aws.String(region),
			},
			SharedConfigState: session.SharedConfigEnable,
			Profile:           c.Profile,
		}))
	}

	if c.HasKeys() {
		return Log(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region: aws.String(region),
				Credentials: credentials.NewStaticCredentials(
					strings.TrimSpace(c.AccessKeyID),
					strings.TrimSpace(c.SecretAccessKey),
					"",
				)}}))
	}

	return nil, fmt.Errorf("You have to specify a profile or credentials for at least one region.")
}

func (c *Credentials) Session(region string) (*session.Session, error) {
	sess, ok := c.cache[region]
	if ok {
		return sess, nil
	}

	sess, err := c.NewSession(region)
	if err != nil {
		c.cache[region] = sess
	}

	return sess, err
}

func Log(s *session.Session, err error) (*session.Session, error) {
	if err == nil {
		s.Handlers.Send.PushFront(func(r *request.Request) {
			log.Debugf("sending AWS request:\n%s", DumpRequest(r.HTTPRequest))
		})

		s.Handlers.ValidateResponse.PushFront(func(r *request.Request) {
			log.Debugf("received AWS response:\n%s", DumpResponse(r.HTTPResponse))
		})
	}

	return s, err
}
