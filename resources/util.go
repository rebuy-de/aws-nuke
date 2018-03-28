package resources

import "github.com/aws/aws-sdk-go/aws/awserr"

func UnPtrBool(ptr *bool, def bool) bool {
	if ptr == nil {
		return def
	}
	return *ptr
}

func UnPtrString(ptr *string, def string) string {
	if ptr == nil {
		return def
	}
	return *ptr
}

func EqualStringPtr(v1, v2 *string) bool {
	if v1 == nil && v2 == nil {
		return true
	}

	if v1 == nil || v2 == nil {
		return false
	}

	return *v1 == *v2
}

func IsAWSError(err error, code string) bool {
	aerr, ok := err.(awserr.Error)
	if !ok {
		return false
	}

	return aerr.Code() == code
}
