package cmd

import (
	"fmt"
	"sort"

	"github.com/rebuy-de/aws-nuke/pkg/awsutil"
	"github.com/rebuy-de/aws-nuke/resources"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	var (
		params NukeParameters
		creds  awsutil.Credentials
	)

	command := &cobra.Command{
		Use:   "aws-nuke",
		Short: "aws-nuke removes every resource from AWS",
		Long:  `A tool which removes every resource from an AWS account.  Use it with caution, since it cannot distinguish between production and non-production.`,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		var err error

		err = params.Validate()
		if err != nil {
			return err
		}

		err = creds.Validate()
		if err != nil {
			return err
		}

		command.SilenceUsage = true

		account, err := awsutil.NewAccount(creds)
		if err != nil {
			return err
		}

		n := NewNuke(params, *account)

		n.Config, err = LoadConfig(n.Parameters.ConfigPath)
		if err != nil {
			return err
		}

		return n.Run()
	}

	command.PersistentFlags().StringVarP(
		&params.ConfigPath, "config", "c", "",
		"path to config (required)")
	command.PersistentFlags().StringVar(
		&creds.Profile, "profile", "",
		"profile name to nuke")
	command.PersistentFlags().StringVar(
		&creds.AccessKeyID, "access-key-id", "",
		"AWS access-key-id")
	command.PersistentFlags().StringVar(
		&creds.SecretAccessKey, "secret-access-key", "",
		"AWS secret-access-key")
	command.PersistentFlags().StringSliceVarP(
		&params.Targets, "target", "t", []string{},
		"limit nuking to certain resource types (eg IAMServerCertificate)")
	command.PersistentFlags().BoolVar(
		&params.NoDryRun, "no-dry-run", false,
		"actually delete found resources")
	command.PersistentFlags().BoolVar(
		&params.Force, "force", false,
		"don't ask for confirmation")

	command.AddCommand(NewVersionCommand())
	command.AddCommand(NewResourceTypesCommand())

	return command
}

func NewResourceTypesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource-types",
		Short: "lists all available resource types",
		Run: func(cmd *cobra.Command, args []string) {
			types := []string{}
			for resourceType, _ := range resources.GetListers() {
				types = append(types, resourceType)
			}

			sort.Strings(types)

			for _, resourceType := range types {
				fmt.Println(resourceType)
			}
		},
	}

	return cmd
}
