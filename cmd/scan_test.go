package cmd

import (
	"strings"
	"testing"

	"github.com/rebuy-de/aws-nuke/resources"
)

func TestSafeLister(t *testing.T) {
	nilLister := func() ([]resources.Resource, error) {
		var ptr *string = nil
		_ = *ptr

		return nil, nil
	}

	_, err := safeLister(nilLister)
	if !strings.Contains(err.Error(), "runtime error: invalid memory address or nil pointer dereference") {
		t.Fatalf("Got unexpected error: %v", err)
	}
}
