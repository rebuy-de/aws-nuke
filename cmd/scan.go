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
		listeners := resources.GetListers(sess)

		for _, lister := range listeners {
			var r []resources.Resource
			r, err = lister()
			if err != nil {
				break
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

	return &Scanner{items, err}
}
