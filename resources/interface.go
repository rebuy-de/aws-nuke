package resources

type ResourceLister func() ([]Resource, error)

type Resource interface {
	Remove() error
	String() string
}

type Waiter interface {
	Resource
	Wait() error
}

type Checker interface {
	Resource
	Check() error
}
