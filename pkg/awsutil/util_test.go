package awsutil_test

import (
	"fmt"
	"testing"

	"github.com/rebuy-de/aws-nuke/v2/pkg/awsutil"
)

func TestSecretRegex(t *testing.T) {
	cases := []struct{ in, out string }{
		{
			in:  "GET / HTTP/1.1\nAuthorization: Never gonna give you up\nHost: bish",
			out: "GET / HTTP/1.1\nAuthorization: <hidden>\nHost: bish",
		},
		{
			in:  "GET / HTTP/1.1\nX-Amz-Security-Token: Never gonna let you down\nHost: bash",
			out: "GET / HTTP/1.1\nX-Amz-Security-Token: <hidden>\nHost: bash",
		},
		{
			in:  "GET / HTTP/1.1\nX-Amz-Security-Token: Never gonna run around and desert you\nAuthorization: Never gonna make you cry",
			out: "GET / HTTP/1.1\nX-Amz-Security-Token: <hidden>\nAuthorization: <hidden>",
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			want := tc.out
			have := string(awsutil.HideSecureHeaders([]byte(tc.in)))

			if want != have {
				t.Errorf("Assertion failed. Want: %#v. Have: %#v", want, have)
			}
		})
	}
}
