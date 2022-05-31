package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rebuy-de/aws-nuke/v2/pkg/config"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ResourceListers map[string]ResourceLister

type ResourceLister func(s *session.Session) ([]Resource, error)

type Resource interface {
	Remove() error
}

type Filter interface {
	Resource
	Filter() error
}

type LegacyStringer interface {
	Resource
	String() string
}

type ResourcePropertyGetter interface {
	Resource
	Properties() types.Properties
}

type FeatureFlagGetter interface {
	Resource
	FeatureFlags(config.FeatureFlags)
}

var resourceListers = make(ResourceListers)

func register(name string, lister ResourceLister, opts ...registerOption) {
	_, exists := resourceListers[name]
	if exists {
		panic(fmt.Sprintf("a resource with the name %s already exists", name))
	}

	resourceListers[name] = lister

	for _, opt := range opts {
		opt(name, lister)
	}
}

var cloudControlMapping = map[string]string{}

func GetCloudControlMapping() map[string]string {
	return cloudControlMapping
}

type registerOption func(name string, lister ResourceLister)

func mapCloudControl(typeName string) registerOption {
	return func(name string, lister ResourceLister) {
		_, exists := cloudControlMapping[typeName]
		if exists {
			panic(fmt.Sprintf("a cloud control mapping for %s already exists", typeName))
		}

		cloudControlMapping[typeName] = name
	}
}

func GetLister(name string) ResourceLister {
	if strings.HasPrefix(name, "AWS::") {
		return NewListCloudControlResource(name)
	}
	return resourceListers[name]
}

func GetListerNames() []string {
	names := []string{}
	for resourceType := range resourceListers {
		names = append(names, resourceType)
	}

	return names
}

func registerCloudControl(typeName string) {
	register(typeName, NewListCloudControlResource(typeName), mapCloudControl(typeName))
}
