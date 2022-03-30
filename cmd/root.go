package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/rebuy-de/aws-nuke/v2/pkg/awsutil"
	"github.com/rebuy-de/aws-nuke/v2/pkg/config"
	"github.com/rebuy-de/aws-nuke/v2/resources"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	var (
		params        NukeParameters
		creds         awsutil.Credentials
		defaultRegion string
		verbose       bool
	)

	command := &cobra.Command{
		Use:   "aws-nuke",
		Short: "aws-nuke removes every resource from AWS",
		Long:  `A tool which removes every resource from an AWS account.  Use it with caution, since it cannot distinguish between production and non-production.`,
	}

	command.PreRun = func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.InfoLevel)
		if verbose {
			log.SetLevel(log.DebugLevel)
		}
		log.SetFormatter(&log.TextFormatter{
			EnvironmentOverrideColors: true,
		})
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		var err error

		err = params.Validate()
		if err != nil {
			return err
		}

		if !creds.HasKeys() && !creds.HasProfile() && defaultRegion != "" {
			creds.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
			creds.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
		}
		err = creds.Validate()
		if err != nil {
			return err
		}

		command.SilenceUsage = true

		config, err := config.Load(params.ConfigPath)
		if err != nil {
			log.Errorf("Failed to parse config file %s", params.ConfigPath)
			return err
		}

		if defaultRegion != "" {
			awsutil.DefaultRegionID = defaultRegion
			switch defaultRegion {
			case endpoints.UsEast1RegionID, endpoints.UsEast2RegionID, endpoints.UsWest1RegionID, endpoints.UsWest2RegionID:
				awsutil.DefaultAWSPartitionID = endpoints.AwsPartitionID
			case endpoints.UsGovEast1RegionID, endpoints.UsGovWest1RegionID:
				awsutil.DefaultAWSPartitionID = endpoints.AwsUsGovPartitionID
			default:
				if config.CustomEndpoints.GetRegion(defaultRegion) == nil {
					err = fmt.Errorf("The custom region '%s' must be specified in the configuration 'endpoints'", defaultRegion)
					log.Error(err.Error())
					return err
				}
			}
		}

		account, err := awsutil.NewAccount(creds, config.CustomEndpoints)
		if err != nil {
			return err
		}

		n := NewNuke(params, *account)

		n.Config = config

		return n.Run()
	}

	command.PersistentFlags().BoolVarP(
		&verbose, "verbose", "v", false,
		"Enables debug output.")

	command.PersistentFlags().StringVarP(
		&params.ConfigPath, "config", "c", "",
		"(required) Path to the nuke config file.")

	command.PersistentFlags().StringVar(
		&creds.Profile, "profile", "",
		"Name of the AWS profile name for accessing the AWS API. "+
			"Cannot be used together with --access-key-id and --secret-access-key.")
	command.PersistentFlags().StringVar(
		&creds.AccessKeyID, "access-key-id", "",
		"AWS access key ID for accessing the AWS API. "+
			"Must be used together with --secret-access-key. "+
			"Cannot be used together with --profile.")
	command.PersistentFlags().StringVar(
		&creds.SecretAccessKey, "secret-access-key", "",
		"AWS secret access key for accessing the AWS API. "+
			"Must be used together with --access-key-id. "+
			"Cannot be used together with --profile.")
	command.PersistentFlags().StringVar(
		&creds.SessionToken, "session-token", "",
		"AWS session token for accessing the AWS API. "+
			"Must be used together with --access-key-id and --secret-access-key. "+
			"Cannot be used together with --profile.")
	command.PersistentFlags().StringVar(
		&creds.AssumeRoleArn, "assume-role-arn", "",
		"AWS IAM role arn to assume. "+
			"The credentials provided via --access-key-id or --profile must "+
			"be allowed to assume this role. ")
	command.PersistentFlags().StringVar(
		&defaultRegion, "default-region", "",
		"Custom default region name.")

	command.PersistentFlags().StringSliceVarP(
		&params.Targets, "target", "t", []string{},
		"Limit nuking to certain resource types (eg IAMServerCertificate). "+
			"This flag can be used multiple times.")
	command.PersistentFlags().StringSliceVarP(
		&params.Excludes, "exclude", "e", []string{},
		"Prevent nuking of certain resource types (eg IAMServerCertificate). "+
			"This flag can be used multiple times.")
	command.PersistentFlags().StringSliceVar(
		&params.CloudControl, "cloud-control", []string{},
		"Nuke given resource via Cloud Control API. "+
			"If there is an old-style method for the same resource, the old-style one will not be executed. "+
			"Note that old-style and cloud-control filters are not compatible! "+
			"This flag can be used multiple times.")
	command.PersistentFlags().BoolVar(
		&params.NoDryRun, "no-dry-run", false,
		"If specified, it actually deletes found resources. "+
			"Otherwise it just lists all candidates.")
	command.PersistentFlags().BoolVar(
		&params.Force, "force", false,
		"Don't ask for confirmation before deleting resources. "+
			"Instead it waits 15s before continuing. Set --force-sleep to change the wait time.")
	command.PersistentFlags().IntVar(
		&params.ForceSleep, "force-sleep", 15,
		"If specified and --force is set, wait this many seconds before deleting resources. "+
			"Defaults to 15.")
	command.PersistentFlags().IntVar(
		&params.MaxWaitRetries, "max-wait-retries", 0,
		"If specified, the program will exit if resources are stuck in waiting for this many iterations. "+
			"0 (default) disables early exit.")
	command.PersistentFlags().BoolVarP(
		&params.Quiet, "quiet", "q", false,
		"Don't show filtered resources.")

	command.AddCommand(NewVersionCommand())
	command.AddCommand(NewResourceTypesCommand())

	return command
}

func NewResourceTypesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource-types",
		Short: "lists all available resource types",
		Run: func(cmd *cobra.Command, args []string) {
			names := resources.GetListerNames()
			sort.Strings(names)

			for _, resourceType := range names {
				fmt.Println(resourceType)
			}
		},
	}

	return cmd
}
