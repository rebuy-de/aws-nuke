package awsutil

import (
	"bytes"
	"net/http"
	"net/http/httputil"
	"regexp"

	"github.com/rebuy-de/aws-nuke/v2/pkg/util"
	log "github.com/sirupsen/logrus"
)

var (
	RESecretHeader = regexp.MustCompile(`(?m:^([^:]*(Auth|Security)[^:]*):.*$)`)
)

func HideSecureHeaders(dump []byte) []byte {
	return RESecretHeader.ReplaceAll(dump, []byte("$1: <hidden>"))
}

func DumpRequest(r *http.Request) string {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.WithField("Error", err).
			Warnf("failed to dump HTTP request")
		return ""
	}

	dump = bytes.TrimSpace(dump)
	dump = HideSecureHeaders(dump)
	dump = util.IndentBytes(dump, []byte("    > "))
	return string(dump)
}

func DumpResponse(r *http.Response) string {
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.WithField("Error", err).
			Warnf("failed to dump HTTP response")
		return ""
	}

	dump = bytes.TrimSpace(dump)
	dump = util.IndentBytes(dump, []byte("    < "))
	return string(dump)
}
