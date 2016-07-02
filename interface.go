package main

type ResourceLister func() ([]Resource, error)

type Resource interface {
	Remove() error
	String() string
}

type Waiter interface {
	Wait() error
}

type Checker interface {
	Check() error
}
