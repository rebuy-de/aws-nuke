package resources

import "github.com/aws/aws-sdk-go/aws/session"

type ResourceListers map[string]ResourceLister

type ResourceLister func(s *session.Session) ([]Resource, error)

type Resource interface {
	Remove() error
	String() string
}

type Filter interface {
	Resource
	Filter() error
}

var resourceListers = make(ResourceListers)

func register(name string, lister ResourceLister) {
	resourceListers[name] = lister
}

func GetListers() ResourceListers {
	return resourceListers
}
