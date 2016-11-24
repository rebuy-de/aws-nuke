package resources

import (
	"fmt"
	"strings"
)

func GetCategory(r Resource) string {
	return strings.Split(fmt.Sprintf("%T", r), ".")[1] // hackey
}
