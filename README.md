# aws-nuke

Nuke a whole AWS account and delete all its resources.

## Caution!

Be aware that *aws-nuke* is a very descructive tool, hence you have to be very
careful while using it. Otherwise you might delete production data.

To reduce the blast radius of accidents, there are some safety precautions:

1. By default *aws-nuke* only lists all nukeable resources. You need to add
   `--no-dry-run` to actually delete resources.
2. *aws-nuke* asks you twice to confirm the deletion by entering the account
   alias. The first time is directly after the start and the second time after
   listing all nukeable resources.
3. To avoid just displaying a account ID, which might gladly be ignored by
   humans, it is required to actually set an [Account
   Alias](http://docs.aws.amazon.com/IAM/latest/UserGuide/console_account-alias.html)
   for your account. Otherwise *aws-nuke* will abort.
4. The Account Alias must not contain the string `prod`. This string is
   hardcoded and it is recommended to add it to every actual production account
   (eg `mycompany-production-ecr`).
5. The config file contains a blacklist field. If the Account ID of the account
   you want to nuke is part of this blacklist, *aws-nuke* will abort. It is
   recommended, that you add every production account to this blacklist.
6. To ensure you just ignore the blacklisting feature, the blacklist must
   contains at least one Account ID.
7. The config file contains account specific settings (eg. filters). The
   account you want to nuke must be explicitly listed there.
8. To ensure to not accidentally delete a random account, it is required to
   specify a config file. It is recommended to have only a single config file
   and add it to a central repository. This way the account blacklist is way
   easier to manage and keep up to date.

Feel free to create an issue, if you have any ideas to improve the safety
procedures.

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
