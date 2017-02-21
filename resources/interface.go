package resources

type ResourceLister func() ([]Resource, error)

type Resource interface {
	Remove() error
	String() string
}

type Filter interface {
	Resource
	Filter() error
}
