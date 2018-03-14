package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rebuy-de/aws-nuke/pkg/awsutil"
	"github.com/rebuy-de/aws-nuke/pkg/util"
	"github.com/rebuy-de/aws-nuke/resources"
	log "github.com/sirupsen/logrus"
)

func Scan(region Region, resourceTypes []string) <-chan *Item {
	items := make(chan *Item, 100)

	go func() {
		for _, resourceType := range resourceTypes {
			lister := resources.GetLister(resourceType)
			rs, err := safeLister(region.Session, lister)
			if err != nil {
				_, ok := err.(awsutil.ErrSkipRequest)
				if ok {
					log.Debug("skipping request: %v", err)
					continue
				}

				dump := util.Indent(fmt.Sprintf("%v", err), "    !!! ")
				log.Errorf("Listing %s failed. Please report this to https://github.com/rebuy-de/aws-nuke/issues/new.\n%s", resourceType, dump)
				continue
			}

			for _, r := range rs {
				items <- &Item{
					Region:   region,
					Resource: r,
					State:    ItemStateNew,
					Type:     resourceType,
				}
			}
		}

		close(items)
	}()

	return items
}

func safeLister(sess *session.Session, lister resources.ResourceLister) (r []resources.Resource, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v\n\n%s", r.(error), string(debug.Stack()))
		}
	}()

	r, err = lister(sess)
	return
}
