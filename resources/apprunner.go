package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apprunner"
	"github.com/rebuy-de/aws-nuke/v2/pkg/awsutil"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AppRunner struct {
	svc         *apprunner.AppRunner
	serviceARN  *string
	serviceName *string
	tags        []*apprunner.Tag
}

func init() {
	register("AppRunner", ListAppRunners)
}

func ListAppRunners(sess *session.Session) ([]Resource, error) {
	svc := apprunner.New(sess)
	resources := make([]Resource, 0)

	params := &apprunner.ListServicesInput{
		MaxResults: aws.Int64(20),
	}

	for {
		resp, err := svc.ListServices(params)
		if err != nil {
			// The ErrUnknownEndpoint occurs when the region doesn't support AppRunner so we will
			// skip those regions
			if _, ok := err.(awsutil.ErrUnknownEndpoint); ok {
				return resources, nil
			}
			return nil, err
		}

		for _, appRunnerService := range resp.ServiceSummaryList {
			tagParams := &apprunner.ListTagsForResourceInput{
				ResourceArn: appRunnerService.ServiceArn,
			}

			tagResp, tagErr := svc.ListTagsForResource(tagParams)
			if tagErr != nil {
				return nil, tagErr
			}

			resources = append(resources, &AppRunner{
				svc:         svc,
				serviceARN:  appRunnerService.ServiceArn,
				serviceName: appRunnerService.ServiceName,
				tags:        tagResp.Tags,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *AppRunner) Remove() error {

	_, err := f.svc.DeleteService(&apprunner.DeleteServiceInput{
		ServiceArn: f.serviceARN,
	})

	return err
}

func (f *AppRunner) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tag := range f.tags {
		properties.SetTag(tag.Key, tag.Value)
	}
	properties.Set("ServiceName", f.serviceName)

	return properties
}

func (f *AppRunner) String() string {
	return *f.serviceARN
}
