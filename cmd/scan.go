package cmd

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rebuy-de/aws-nuke/resources"
)

type Scanner struct {
	Items <-chan *Item
	Error error
}

func Scan(sess *session.Session) *Scanner {
	var err error
	items := make(chan *Item, 100)

	go func() {
		listers := resources.GetListers(sess)

		for _, lister := range listers {
			var r []resources.Resource
			r, err = lister()
			if err != nil {
				break
			}

			for _, r := range r {
				items <- &Item{
					Resource: r,
					Service:  resources.GetCategory(r),
					Lister:   lister,
					State:    ItemStateNew,
				}
			}
		}

		close(items)
	}()

	return &Scanner{items, err}
}
