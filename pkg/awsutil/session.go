package awsutil

import (
	"fmt"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
)

const (
	GlobalRegionID  = "global"
	DefaultRegionID = endpoints.UsEast1RegionID
)

type Credentials struct {
	Profile string

	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string

	session *session.Session
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

func (c *Credentials) rootSession() (*session.Session, error) {
	if c.session == nil {
		var opts session.Options

		region := DefaultRegionID
		log.Debugf("creating new root session in %s", region)

		switch {
		case c.HasProfile() == c.HasKeys():
			return nil, fmt.Errorf("You have to specify a profile or credentials for at least one region.")

		case c.HasProfile():
			opts = session.Options{
				SharedConfigState:       session.SharedConfigEnable,
				Profile:                 c.Profile,
				AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
			}

		case c.HasKeys():
			opts = session.Options{
				Config: aws.Config{
					Credentials: credentials.NewStaticCredentials(
						strings.TrimSpace(c.AccessKeyID),
						strings.TrimSpace(c.SecretAccessKey),
						strings.TrimSpace(c.SessionToken),
					)}}
		}

		opts.Config.Region = aws.String(region)
		opts.Config.DisableRestProtocolURICleaning = aws.Bool(true)

		sess, err := session.NewSessionWithOptions(opts)
		if err != nil {
			return nil, err
		}

		c.session = sess
	}

	return c.session, nil
}

func (c *Credentials) NewSession(region string) (*session.Session, error) {
	log.Debugf("creating new session in %s", region)

	global := false

	if region == GlobalRegionID {
		region = DefaultRegionID
		global = true
	}

	root, err := c.rootSession()
	if err != nil {
		return nil, err
	}

	sess := root.Copy(&aws.Config{
		Region: &region,
	})

	sess.Handlers.Send.PushFront(func(r *request.Request) {
		log.Debugf("sending AWS request:\n%s", DumpRequest(r.HTTPRequest))
	})

	sess.Handlers.ValidateResponse.PushFront(func(r *request.Request) {
		log.Debugf("received AWS response:\n%s", DumpResponse(r.HTTPResponse))
	})

	sess.Handlers.Validate.PushFront(skipMissingServiceInRegionHandler)
	sess.Handlers.Validate.PushFront(skipGlobalHandler(global))

	return sess, nil
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

func skipGlobalHandler(global bool) func(r *request.Request) {
	return func(r *request.Request) {
		service := r.ClientInfo.ServiceName

		rs, ok := endpoints.RegionsForService(endpoints.DefaultPartitions(), endpoints.AwsPartitionID, service)
		if !ok {
			// This means that the service does not exist in the endpoints list.
			if global {
				r.Error = ErrSkipRequest(fmt.Sprintf("service '%s' is was not found in the endpoint list; assuming it is not global", service))
			} else {
				host := r.HTTPRequest.URL.Hostname()
				_, err := net.LookupHost(host)
				if err != nil {
					log.Debug(err)
					r.Error = ErrUnknownEndpoint(fmt.Sprintf("DNS lookup failed for %s; assuming it does not exist in this region", host))
				}
			}
			return
		}

		if len(rs) == 0 && !global {
			r.Error = ErrSkipRequest(fmt.Sprintf("service '%s' is global, but the session is not", service))
			return
		}

		if len(rs) > 0 && global {
			r.Error = ErrSkipRequest(fmt.Sprintf("service '%s' is not global, but the session is", service))
			return
		}
	}
}
