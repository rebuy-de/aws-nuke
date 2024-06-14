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

func Chunk[T any](slice []T, size int) [][]T {
	var chunks [][]T
	for i := 0; i < len(slice); {
		// Clamp the last chunk to the slice bound as necessary.
		end := size
		if l := len(slice[i:]); l < size {
			end = l
		}

		// Set the capacity of each chunk so that appending to a chunk does not
		// modify the original slice.
		chunks = append(chunks, slice[i:i+end:i+end])
		i += end
	}

	return chunks
}
