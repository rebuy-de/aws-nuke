package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rebuy-de/aws-nuke/pkg/types"
)

func Prompt(expect string) error {
	fmt.Print("> ")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	if strings.TrimSpace(text) != expect {
		return fmt.Errorf("aborted")
	}
	fmt.Println()

	return nil
}

func ResolveResourceTypes(base types.Collection, include, exclude []types.Collection) types.Collection {
	for _, i := range include {
		if len(i) > 0 {
			base = base.Intersect(i)
		}
	}

	for _, e := range exclude {
		base = base.Remove(e)
	}

	return base
}
