package cmd

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rebuy-de/aws-nuke/resources"
)

func TestSafeLister(t *testing.T) {
	nilLister := func(s *session.Session) ([]resources.Resource, error) {
		// generate nil pointer dereference panic
		var ptr *string
		_ = *ptr

		return nil, nil
	}

	_, err := safeLister(nil, nilLister)
	if !strings.Contains(err.Error(), "runtime error: invalid memory address or nil pointer dereference") {
		t.Fatalf("Got unexpected error: %v", err)
	}
}
