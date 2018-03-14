package awsutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
)

type Credentials struct {
	Profile string

	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string

	cache map[string]*session.Session
}

func (c *Credentials) HasProfile() bool {
	return strings.TrimSpace(c.Profile) != ""
}

func (c *Credentials) HasKeys() bool {
	return strings.TrimSpace(c.AccessKeyID) != "" ||
		strings.TrimSpace(c.SecretAccessKey) != "" ||
		strings.TrimSpace(c.SessionToken) != ""
}

func (c *Credentials) Validate() error {
	if c.HasProfile() == c.HasKeys() {
		return fmt.Errorf("You have to specify the --profile flag OR " +
			"--access-key-id with --secret-access-key and optionally " +
			"--session-token.\n")
	}

	return nil
}

func (c *Credentials) NewSession(region string) (*session.Session, error) {
	var opts session.Options

	switch {
	case c.HasProfile() && c.HasKeys():
		return nil, fmt.Errorf("You have to specify a profile or credentials for at least one region.")

	case c.HasProfile():
		opts = session.Options{
			Config: aws.Config{
				Region: aws.String(region),
			},
			SharedConfigState: session.SharedConfigEnable,
			Profile:           c.Profile,
		}

	case c.HasKeys():
		opts = session.Options{
			Config: aws.Config{
				Region: aws.String(region),
				Credentials: credentials.NewStaticCredentials(
					strings.TrimSpace(c.AccessKeyID),
					strings.TrimSpace(c.SecretAccessKey),
					strings.TrimSpace(c.SessionToken),
				)}}
	}

	sess, err := session.NewSessionWithOptions(opts)
	if err != nil {
		return nil, err
	}

	addHandlers(sess)

	return sess, nil
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

func addHandlers(s *session.Session) {
	s.Handlers.Validate.PushFront(skipMissingServiceInRegionHandler)

	s.Handlers.Send.PushFront(func(r *request.Request) {
		log.Debugf("sending AWS request:\n%s", DumpRequest(r.HTTPRequest))
	})

	s.Handlers.ValidateResponse.PushFront(func(r *request.Request) {
		log.Debugf("received AWS response:\n%s", DumpResponse(r.HTTPResponse))
	})
}

func skipMissingServiceInRegionHandler(r *request.Request) {
	region := *r.Config.Region
	service := r.ClientInfo.ServiceName

	rs, ok := endpoints.RegionsForService(endpoints.DefaultPartitions(), endpoints.AwsPartitionID, service)
	if !ok {
		// This means that the service does not exist and this shouldn't be handled here.
		return
	}

	if len(rs) == 0 {
		// Avoid to throw an error on global services.
		return
	}

	_, ok = rs[region]
	if !ok {
		r.Error = ErrSkipRequest(fmt.Sprintf(
			"service '%s' is not available in region '%s'",
			service, region))
	}
}
