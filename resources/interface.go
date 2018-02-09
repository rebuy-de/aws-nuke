package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
)

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
	_, exists := resourceListers[name]
	if exists {
		panic(fmt.Sprintf("a resource with the name %s already exists", name))
	}

	resourceListers[name] = lister
}

func GetListers() ResourceListers {
	return resourceListers
}

func GetLister(name string) ResourceLister {
	return resourceListers[name]
}

func GetListerNames() []string {
	names := []string{}
	for resourceType, _ := range GetListers() {
		names = append(names, resourceType)
	}

	return names
}
