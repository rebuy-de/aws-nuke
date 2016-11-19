package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/rebuy-de/aws-nuke/resources"
)

var (
	ReasonSkip            = *color.New(color.FgYellow)
	ReasonError           = *color.New(color.FgRed)
	ReasonRemoveTriggered = *color.New(color.FgGreen)
	ReasonWaitPending     = *color.New()
	ReasonSuccess         = *color.New(color.FgGreen)
	ColorID               = *color.New(color.Bold)
)

func Log(r resources.Resource, c color.Color, msg string) {
	fmt.Printf("[%s] ", time.Now().Format(time.RFC3339))
	fmt.Print(strings.Split(fmt.Sprintf("%T", r), ".")[1]) // hackey
	fmt.Printf(" - ")
	ColorID.Printf("'%s'", r.String())
	fmt.Printf(" - ")
	c.Printf("%s\n", msg)
}

func LogErrorf(err error) {
	out := color.New(color.FgRed)
	trace := fmt.Sprintf("%+v", err)
	out.Println(trace)
	out.Println("")
}
