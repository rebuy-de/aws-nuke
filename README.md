# aws-nuke
Nuke a whole AWS account

## Usage

```
Usage:
  aws-nuke [flags]
  aws-nuke [command]

Available Commands:
  version     shows version of this application

Flags:
      --access-key-id string       AWS access-key-id
  -c, --config string              path to config (required)
      --force                      don't ask for confirmation
      --no-dry-run                 actualy delete found resources
      --profile string             profile name to nuke
      --secret-access-key string   AWS secret-access-key
  -t, --target stringSlice         limit nuking to certain resource types (eg IamServerCertificate)

Use "aws-nuke [command] --help" for more information about a command.
```
