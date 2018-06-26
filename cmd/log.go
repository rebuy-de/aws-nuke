package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rebuy-de/aws-nuke/resources"
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

func Log(region Region, resourceType string, r resources.Resource, c color.Color, msg string) {
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
		ColorResourceProperties.Print(rProp.Properties())
		fmt.Printf(" - ")
	}

	c.Printf("%s\n", msg)
}
