package cmd

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudcontrolapi"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/rebuy-de/aws-nuke/pkg/types"
	"github.com/rebuy-de/aws-nuke/pkg/util"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func CloudControlGetListerNames(ctx context.Context, sess *session.Session) ([]string, error) {
	cf := cloudformation.New(sess)

	in := &cloudformation.ListTypesInput{
		Type:             aws.String(cloudformation.RegistryTypeResource),
		Visibility:       aws.String(cloudformation.VisibilityPublic),
		ProvisioningType: aws.String(cloudformation.ProvisioningTypeFullyMutable),
	}

	names := []string{}
	err := cf.ListTypesPagesWithContext(ctx, in, func(out *cloudformation.ListTypesOutput, _ bool) bool {
		if out == nil {
			return true
		}

		for _, summary := range out.TypeSummaries {
			if summary == nil {
				continue
			}

			name := aws.StringValue(summary.TypeName)
			logrus.Debugf("Registering %s", name)
			names = append(names, name)
		}

		return true
	})
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Registered %d resources", len(names))

	return names, nil
}

type CloudControlScanner struct {
	items chan *Item
	sess  *session.Session
	ctx   context.Context
}

func CloudControlScan(ctx context.Context, sess *session.Session, region *Region, resourceTypes []string) <-chan *Item {
	s := &CloudControlScanner{
		items: make(chan *Item, 100),
		sess:  sess,
		ctx:   ctx,
	}
	go s.run(region, resourceTypes)

	return s.items
}

func (s *CloudControlScanner) run(region *Region, resourceTypes []string) {
	for _, resourceType := range resourceTypes {
		s.list(region, resourceType)
	}

	close(s.items)
}

func (s *CloudControlScanner) list(region *Region, resourceType string) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v\n\n%s", r.(error), string(debug.Stack()))
			dump := util.Indent(fmt.Sprintf("%v", err), "    ")
			log.Errorf("Listing %s failed:\n%s", resourceType, dump)
		}
	}()

	cc := cloudcontrolapi.New(s.sess)
	in := &cloudcontrolapi.ListResourcesInput{
		TypeName: aws.String(resourceType),
	}
	err := cc.ListResourcesPagesWithContext(s.ctx, in, func(page *cloudcontrolapi.ListResourcesOutput, lastPage bool) bool {
		for _, desc := range page.ResourceDescriptions {
			s.items <- &Item{
				Region: region,
				State:  ItemStateNew,
				Type:   resourceType,
				Resource: &CloudControlResource{
					identifier: aws.StringValue(desc.Identifier),
					properties: aws.StringValue(desc.Properties),
				},
			}
		}

		return true
	})
	if err != nil {
		logrus.WithError(err).Errorf("failed to list %s", resourceType)
	}

	time.Sleep(time.Second)
}

type CloudControlResource struct {
	identifier string
	properties string
}

func (r *CloudControlResource) String() string {
	return r.identifier
}

func (r *CloudControlResource) Properties() types.Properties {
	return types.NewProperties().
		Set("Identifier", r.identifier).
		Set("Properties", r.properties)
}

func (f *CloudControlResource) Remove() error {
	return nil
}
