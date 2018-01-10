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
	ColorRegion       = *color.New(color.Bold)
	ColorResourceType = *color.New()
	ColorResourceID   = *color.New(color.Bold)
	ColorWarning      = *color.New(color.FgYellow)
)

func Log(region Region, resourceType string, r resources.Resource, c color.Color, msg string) {
	ColorRegion.Printf("%s", region.Name)
	fmt.Printf(" - ")
	ColorResourceType.Print(resourceType)
	fmt.Printf(" - ")
	ColorResourceID.Printf("'%s'", r.String())
	fmt.Printf(" - ")
	c.Printf("%s\n", msg)
}

func LogWarn(s string, i ...interface{}) {
	ColorWarning.Printf("WARNING: "+s, i...)
}

func LogErrorf(err error) {
	out := color.New(color.FgRed)
	trace := fmt.Sprintf("%+v", err)
	out.Println(trace)
	out.Println("")
}
