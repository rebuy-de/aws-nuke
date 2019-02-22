package cmd

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/rebuy-de/aws-nuke/pkg/awsutil"
	"github.com/rebuy-de/aws-nuke/pkg/util"
	"github.com/rebuy-de/aws-nuke/resources"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

const ScannerParallelQueries = 16

func Scan(region Region, resourceTypes []string) <-chan *Item {
	s := &scanner{
		items:     make(chan *Item, 100),
		semaphore: semaphore.NewWeighted(ScannerParallelQueries),
	}
	go s.run(region, resourceTypes)

	return s.items
}

type scanner struct {
	items     chan *Item
	semaphore *semaphore.Weighted
}

func (s *scanner) run(region Region, resourceTypes []string) {
	ctx := context.Background()

	for _, resourceType := range resourceTypes {
		s.semaphore.Acquire(ctx, 1)
		go s.list(region, resourceType)
	}

	// Wait for all routines to finish.
	s.semaphore.Acquire(ctx, ScannerParallelQueries)

	close(s.items)
}

func (s *scanner) list(region Region, resourceType string) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v\n\n%s", r.(error), string(debug.Stack()))
			dump := util.Indent(fmt.Sprintf("%v", err), "    ")
			log.Errorf("Listing %s failed:\n%s", resourceType, dump)
		}
	}()
	defer s.semaphore.Release(1)

	lister := resources.GetLister(resourceType)
	rs, err := lister(region.Session)
	if err != nil {
		_, ok := err.(awsutil.ErrSkipRequest)
		if ok {
			log.Debugf("skipping request: %v", err)
			return
		}

		_, ok = err.(awsutil.ErrUnknownEndpoint)
		if ok {
			log.Warnf("skipping request: %v", err)
			return
		}

		dump := util.Indent(fmt.Sprintf("%v", err), "    ")
		log.Errorf("Listing %s failed:\n%s", resourceType, dump)
		return
	}

	for _, r := range rs {
		s.items <- &Item{
			Region:   region,
			Resource: r,
			State:    ItemStateNew,
			Type:     resourceType,
		}
	}
}
