package main

import "log"

var (
	// will be overwritten on build
	version = "unknown"
)

func main() {
	log.Printf("aws-nuke %s", version)
}
