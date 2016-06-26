package main

type ResourceLister func() ([]Resource, error)

type Resource interface {
	Check() error
	Remove() error
	Wait() error
	String() string
}

type Skipped struct {
	Reason string
}
