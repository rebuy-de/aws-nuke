package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rebuy-de/aws-nuke/resources"
)

func Scan(sess *session.Session) <-chan *Item {
	items := make(chan *Item, 100)

	go func() {
		listers := resources.GetListers(sess)

		for _, lister := range listers {
			r, err := safeLister(lister)
			if err != nil {
				LogErrorf(fmt.Errorf("\n=============\n\n"+
					"Listing with %T failed:\n\n"+
					"%v\n\n"+
					"Please report this to https://github.com/rebuy-de/aws-nuke/issues/new.\n\n"+
					"=============",
					lister, err))
				continue
			}

			for _, r := range r {
				items <- &Item{
					Region:   *sess.Config.Region,
					Resource: r,
					Service:  resources.GetCategory(r),
					Lister:   lister,
					State:    ItemStateNew,
				}
			}
		}

		close(items)
	}()

	return items
}

func safeLister(lister resources.ResourceLister) (r []resources.Resource, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v\n\n%s", r.(error), string(debug.Stack()))
		}
	}()

	r, err = lister()
	return
}
