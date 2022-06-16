package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/rebuy-de/aws-nuke/v2/resources"
)

var (
	ReasonSkip            = *color.New(color.FgYellow)
	ReasonError           = *color.New(color.FgRed)
	ReasonRemoveTriggered = *color.New(color.FgGreen)
	ReasonWaitPending     = *color.New(color.FgBlue)
	ReasonSuccess         = *color.New(color.FgGreen)
)

var (
	ColorRegion             = *color.New(color.Bold)
	ColorResourceType       = *color.New()
	ColorResourceID         = *color.New(color.Bold)
	ColorResourceProperties = *color.New(color.Italic)
)

// Format the resource properties in sorted order ready for printing.
// This ensures that multiple runs of aws-nuke produce stable output so
// that they can be compared with each other.
func Sorted(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sorted := make([]string, 0, len(m))
	for k := range keys {
		sorted = append(sorted, fmt.Sprintf("%s: \"%s\"", keys[k], m[keys[k]]))
	}
	return fmt.Sprintf("[%s]", strings.Join(sorted, ", "))
}

func Log(region *Region, resourceType string, r resources.Resource, c color.Color, msg string) {
	ColorRegion.Printf("%s", region.Name)
	fmt.Printf(" - ")
	ColorResourceType.Print(resourceType)
	fmt.Printf(" - ")

	rString, ok := r.(resources.LegacyStringer)
	if ok {
		ColorResourceID.Print(rString.String())
		fmt.Printf(" - ")
	}

	rProp, ok := r.(resources.ResourcePropertyGetter)
	if ok {
		ColorResourceProperties.Print(Sorted(rProp.Properties()))
		fmt.Printf(" - ")
	}

	c.Printf("%s\n", msg)
}
