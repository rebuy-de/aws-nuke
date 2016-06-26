package main

type ResourceLister func() ([]Resource, error)

type Resource interface {
	Remove() error
	Wait() error
	String() string
}
