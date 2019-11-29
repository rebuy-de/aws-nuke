package awsutil

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rebuy-de/aws-nuke/pkg/config"
	log "github.com/sirupsen/logrus"
)

const (
	GlobalRegionID = "global"
)

var (
	// DefaultRegionID The default region. Can be customized for non AWS implementations
	DefaultRegionID = endpoints.UsEast1RegionID
)

type Credentials struct {
	Profile string

	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string

	CustomEndpoints config.CustomEndpoints
	session         *session.Session
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
					Credentials: c.awsNewStaticCredentials(),
				},
			}
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

func (c *Credentials) awsNewStaticCredentials() *credentials.Credentials {
	if !c.HasKeys() {
		return credentials.NewEnvCredentials()
	}
	return credentials.NewStaticCredentials(
		strings.TrimSpace(c.AccessKeyID),
		strings.TrimSpace(c.SecretAccessKey),
		strings.TrimSpace(c.SessionToken),
	)
}

func (c *Credentials) NewSession(region, serviceType string) (*session.Session, error) {
	log.Debugf("creating new session in %s for %s", region, serviceType)

	global := false

	if region == GlobalRegionID {
		region = DefaultRegionID
		global = true
	}

	var sess *session.Session
	isCustom := false
	if customRegion := c.CustomEndpoints.GetRegion(region); customRegion != nil {
		customService := customRegion.Services.GetService(serviceType)
		if customService == nil {
			return nil, ErrSkipRequest(fmt.Sprintf(
				".service '%s' is not available in region '%s'",
				serviceType, region))
		}
		conf := &aws.Config{
			Region:      &region,
			Endpoint:    &customService.URL,
			Credentials: c.awsNewStaticCredentials(),
		}
		if customService.TLSInsecureSkipVerify {
			conf.HTTPClient = &http.Client{Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}}
		}
		// ll := aws.LogDebugWithEventStreamBody
		// conf.LogLevel = &ll
		var err error
		sess, err = session.NewSession(conf)
		if err != nil {
			return nil, err
		}
		isCustom = true
	}

	if sess == nil {
		root, err := c.rootSession()
		if err != nil {
			return nil, err
		}

		sess = root.Copy(&aws.Config{
			Region: &region,
		})
	}

	sess.Handlers.Send.PushFront(func(r *request.Request) {
		log.Debugf("sending AWS request:\n%s", DumpRequest(r.HTTPRequest))
	})

	sess.Handlers.ValidateResponse.PushFront(func(r *request.Request) {
		log.Debugf("received AWS response:\n%s", DumpResponse(r.HTTPResponse))
	})

	if !isCustom {
		sess.Handlers.Validate.PushFront(skipMissingServiceInRegionHandler)
		sess.Handlers.Validate.PushFront(skipGlobalHandler(global))
	}
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
