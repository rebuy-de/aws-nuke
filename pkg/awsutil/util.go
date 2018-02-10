package awsutil

import (
	"bytes"
	"net/http"
	"net/http/httputil"
	"regexp"

	log "github.com/sirupsen/logrus"
)

func Indent(s, prefix string) string {
	return string(IndentBytes([]byte(s), []byte(prefix)))
}

func IndentBytes(b, prefix []byte) []byte {
	var res []byte
	bol := true
	for _, c := range b {
		if bol && c != '\n' {
			res = append(res, prefix...)
		}
		res = append(res, c)
		bol = c == '\n'
	}
	return res
}

var REAuthHeader = regexp.MustCompile(`(?m:^(Auth[^:]*):.*$)`)

func DumpRequest(r *http.Request) string {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.WithField("Error", err).
			Warnf("failed to dump HTTP request")
		return ""
	}

	dump = bytes.TrimSpace(dump)
	dump = REAuthHeader.ReplaceAll(dump, []byte("$1: <hidden>"))
	dump = IndentBytes(dump, []byte("    > "))
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
	dump = IndentBytes(dump, []byte("    < "))
	return string(dump)
}
