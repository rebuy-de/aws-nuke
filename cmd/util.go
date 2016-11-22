package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func AskContinue(message string, opts ...interface{}) error {
	fmt.Println()
	fmt.Printf(message, opts...)
	fmt.Println()
	fmt.Println("Do you want to continue? (type 'yes')")

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	if strings.TrimSpace(text) != "yes" {
		return fmt.Errorf("aborted")
	}
	fmt.Println()

	return nil
}
