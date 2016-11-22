package cmd

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	params := NukeParameters{}

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

		command.SilenceUsage = true

		n := NewNuke(params)

		n.Config, err = LoadConfig(n.Parameters.ConfigPath)
		if err != nil {
			return err
		}

		err = n.StartSession()
		if err != nil {
			return err
		}

		return n.Run()
	}

	command.PersistentFlags().StringVarP(
		&params.ConfigPath, "config", "c", "",
		"path to config (required)")
	command.PersistentFlags().StringVar(
		&params.Profile, "profile", "",
		"profile name to nuke")
	command.PersistentFlags().StringVar(
		&params.AccessKeyID, "access-key-id", "",
		"AWS access-key-id")
	command.PersistentFlags().StringVar(
		&params.SecretAccessKey, "secret-access-key", "",
		"AWS secret-access-key")
	command.PersistentFlags().BoolVar(
		&params.NoDryRun, "no-dry-run", false,
		"actualy delete found resources")
	command.PersistentFlags().BoolVar(
		&params.Force, "force", false,
		"don't ask for confirmation")

	return command
}
