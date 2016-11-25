package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
