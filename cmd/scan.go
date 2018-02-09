package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rebuy-de/aws-nuke/resources"
)

func Scan(region Region, resourceTypes []string) <-chan *Item {
	items := make(chan *Item, 100)

	go func() {
		for _, resourceType := range resourceTypes {
			lister := resources.GetLister(resourceType)
			rs, err := safeLister(region.Session, lister)
			if err != nil {
				LogErrorf(fmt.Errorf("\n=============\n\n"+
					"Listing with %T failed:\n\n"+
					"%v\n\n"+
					"Please report this to https://github.com/rebuy-de/aws-nuke/issues/new.\n\n"+
					"=============",
					lister, err))
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
