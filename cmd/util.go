package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
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

func ResolveResourceTypes(
	base types.Collection, mapping map[string]string,
	include, exclude, cloudControl []types.Collection) types.Collection {

	for _, cl := range cloudControl {
		oldStyle := types.Collection{}
		for _, c := range cl {
			os, found := mapping[c]
			if found {
				oldStyle = append(oldStyle, os)
			}
		}

		base = base.Union(cl)
		base = base.Remove(oldStyle)
	}

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

func IsTrue(s string) bool {
	return strings.TrimSpace(strings.ToLower(s)) == "true"
}
