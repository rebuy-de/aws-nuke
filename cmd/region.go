package cmd

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rebuy-de/aws-nuke/v2/pkg/awsutil"
)

// SessionFactory support for custom endpoints
type SessionFactory func(regionName, svcType string) (*session.Session, error)

// ResourceTypeResolver returns the service type from the resourceType
type ResourceTypeResolver func(regionName, resourceType string) string

type Region struct {
	Name            string
	NewSession      SessionFactory
	ResTypeResolver ResourceTypeResolver

	cache map[string]*session.Session
	lock  *sync.RWMutex
}

func NewRegion(name string, typeResolver ResourceTypeResolver, sessionFactory SessionFactory) *Region {
	return &Region{
		Name:            name,
		NewSession:      sessionFactory,
		ResTypeResolver: typeResolver,
		lock:            &sync.RWMutex{},
		cache:           make(map[string]*session.Session),
	}
}

func (region *Region) Session(resourceType string) (*session.Session, error) {
	svcType := region.ResTypeResolver(region.Name, resourceType)
	if svcType == "" {
		return nil, awsutil.ErrSkipRequest(fmt.Sprintf(
			"No service available in region '%s' to handle '%s'",
			region.Name, resourceType))
	}

	// Need to read
	region.lock.RLock()
	sess := region.cache[svcType]
	region.lock.RUnlock()
	if sess != nil {
		return sess, nil
	}

	// Need to write:
	region.lock.Lock()
	sess, err := region.NewSession(region.Name, svcType)
	if err != nil {
		region.lock.Unlock()
		return nil, err
	}
	region.cache[svcType] = sess
	region.lock.Unlock()
	return sess, nil
}
