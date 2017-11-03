package resources

import (
	"fmt"
	"strings"
)

func GetCategory(r Resource) string {
	return strings.Split(fmt.Sprintf("%T", r), ".")[1] // hackey
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
