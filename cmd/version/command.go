package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev" // should be set via ldflags
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Shows current version of aws-nuke",
		Run: func(cmd *cobra.Command, args []string) {
			Print()
		},
	}
}

func Print() {
	fmt.Printf("aws-nuke version %s\n", version)
}
