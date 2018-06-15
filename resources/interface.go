package resources

import (
	"fmt"
	"strings"

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

type Properties map[string]string

func NewProperties() Properties {
	return make(Properties)
}

func (p Properties) String() string {
	parts := []string{}
	for k, v := range p {
		parts = append(parts, fmt.Sprintf(`%s: "%v"`, k, v))
	}

	return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
}

func (p Properties) Set(key string, value interface{}) Properties {
	if value == nil {
		return p
	}

	switch v := value.(type) {
	case *string:
		p[key] = *v
	default:
		// Fallback to Stringer interface. This produces gibberish on pointers,
		// but is the only way to avoid reflection.
		p[key] = fmt.Sprint(value)
	}

	return p
}

func (p Properties) Get(key string) string {
	value, ok := p[key]
	if !ok {
		return ""
	}

	return value
}

type ResourcePropertyGetter interface {
	Resource
	Properties() Properties
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
