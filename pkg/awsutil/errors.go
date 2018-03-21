package awsutil

type ErrSkipRequest string

func (err ErrSkipRequest) Error() string {
	return string(err)
}
