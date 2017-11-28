# aws-nuke

[![Build Status](https://travis-ci.org/rebuy-de/aws-nuke.svg?branch=master)](https://travis-ci.org/rebuy-de/aws-nuke)
[![license](https://img.shields.io/github/license/rebuy-de/aws-nuke.svg)]()
[![GitHub release](https://img.shields.io/github/release/rebuy-de/aws-nuke.svg)]()

Nuke a whole AWS account and delete all its resources.

> **Development Status** *aws-nuke* is stable, but currently not all AWS
resources are covered by it. Be encouraged to add missing resources and create
a Pull Request or to create an [Issue](https://github.com/rebuy-de/aws-nuke/issues/new).

## Caution!

Be aware that *aws-nuke* is a very destructive tool, hence you have to be very
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


## Use Cases

* We are testing our [Terraform](https://www.terraform.io/) code with Jenkins.
  Sometimes a Terraform run fails during development and messes up the account.
  With *aws-nuke* we can simply clean up the failed account so it can be reused
  for the next build.
* Our platform developers have their own AWS Accounts where they can create
  their own Kubernetes clusters for testing purposes. With *aws-nuke* it is
  very easy to clean up these account at the end of the day and keep the costs
  low.


## Usage

At first you need to create a config file for *aws-nuke*. This is a minimal one:

```yaml
regions:
- eu-west-1

account-blacklist:
- "999999999999" # production

accounts:
  "000000000000": {} # aws-nuke-example
```

With this config we can run *aws-nuke*:

```
$ aws-nuke -c config/nuke-config.yml --profile aws-nuke-example
aws-nuke version v1.0.39.gc2f318f - Fri Jul 28 16:26:41 CEST 2017 - c2f318f37b7d2dec0e646da3d4d05ab5296d5bce

Do you really want to nuke the account with the ID 000000000000 and the alias 'aws-nuke-example'?
Do you want to continue? Enter account alias to continue.
> aws-nuke-example

eu-west-1 - EC2DhcpOption - 'dopt-bf2ec3d8' - would remove
eu-west-1 - EC2Instance - 'i-01b489457a60298dd' - would remove
eu-west-1 - EC2KeyPair - 'test' - would remove
eu-west-1 - EC2NetworkACL - 'acl-6482a303' - cannot delete default VPC
eu-west-1 - EC2RouteTable - 'rtb-ffe91e99' - would remove
eu-west-1 - EC2SecurityGroup - 'sg-220e945a' - cannot delete group 'default'
eu-west-1 - EC2SecurityGroup - 'sg-f20f958a' - would remove
eu-west-1 - EC2Subnet - 'subnet-154d844e' - would remove
eu-west-1 - EC2Volume - 'vol-0ddfb15461a00c3e2' - would remove
eu-west-1 - EC2Vpc - 'vpc-c6159fa1' - would remove
eu-west-1 - IAMUserAccessKey - 'my-user -> ABCDEFGHIJKLMNOPQRST' - would remove
eu-west-1 - IAMUserPolicyAttachment - 'my-user -> AdministratorAccess' - would remove
eu-west-1 - IAMUser - 'my-user' - would remove
Scan complete: 13 total, 11 nukeable, 2 filtered.

Would delete these resources. Provide --no-dry-run to actually destroy resources.
```

As we see, *aws-nuke* only lists all found resources and exits. This is because
the `--no-dry-run` flag is missing. Also it wants to delete the
administrator. We don't want to do this, because we use this user to access
our account. Therefore we have to extend the config so it ignores this user:


```yaml
regions:
- eu-west-1

account-blacklist:
- "999999999999" # production

accounts:
  "000000000000": # aws-nuke-example
    filters:
      IAMUser:
      - "my-user"
      IAMUserPolicyAttachment:
      - "my-user -> AdministratorAccess"
      IAMUserAccessKey:
      - "my-user -> ABCDEFGHIJKLMNOPQRST"
```

```
$ aws-nuke -c config/nuke-config.yml --profile aws-nuke-example --no-dry-run
aws-nuke version v1.0.39.gc2f318f - Fri Jul 28 16:26:41 CEST 2017 - c2f318f37b7d2dec0e646da3d4d05ab5296d5bce

Do you really want to nuke the account with the ID 000000000000 and the alias 'aws-nuke-example'?
Do you want to continue? Enter account alias to continue.
> aws-nuke-example

eu-west-1 - EC2DhcpOption - 'dopt-bf2ec3d8' - would remove
eu-west-1 - EC2Instance - 'i-01b489457a60298dd' - would remove
eu-west-1 - EC2KeyPair - 'test' - would remove
eu-west-1 - EC2NetworkACL - 'acl-6482a303' - cannot delete default VPC
eu-west-1 - EC2RouteTable - 'rtb-ffe91e99' - would remove
eu-west-1 - EC2SecurityGroup - 'sg-220e945a' - cannot delete group 'default'
eu-west-1 - EC2SecurityGroup - 'sg-f20f958a' - would remove
eu-west-1 - EC2Subnet - 'subnet-154d844e' - would remove
eu-west-1 - EC2Volume - 'vol-0ddfb15461a00c3e2' - would remove
eu-west-1 - EC2Vpc - 'vpc-c6159fa1' - would remove
eu-west-1 - IAMUserAccessKey - 'my-user -> ABCDEFGHIJKLMNOPQRST' - filtered by config
eu-west-1 - IAMUserPolicyAttachment - 'my-user -> AdministratorAccess' - filtered by config
eu-west-1 - IAMUser - 'my-user' - filtered by config
Scan complete: 13 total, 8 nukeable, 5 filtered.

Do you really want to nuke these resources on the account with the ID 000000000000 and the alias 'aws-nuke-example'?
Do you want to continue? Enter account alias to continue.
> aws-nuke-example

eu-west-1 - EC2DhcpOption - 'dopt-bf2ec3d8' - DependencyViolation: The dhcpOptions 'dopt-bf2ec3d8' has dependencies and cannot be deleted.
	status code: 400, request id: 9665c066-6bb1-4643-9071-f03481f80d4e
eu-west-1 - EC2Instance - 'i-01b489457a60298dd' - triggered remove
eu-west-1 - EC2KeyPair - 'test' - triggered remove
eu-west-1 - EC2RouteTable - 'rtb-ffe91e99' - DependencyViolation: The routeTable 'rtb-ffe91e99' has dependencies and cannot be deleted.
	status code: 400, request id: 3f667620-3207-4576-ae68-0b75261c0504
eu-west-1 - EC2SecurityGroup - 'sg-f20f958a' - DependencyViolation: resource sg-f20f958a has a dependent object
	status code: 400, request id: 5da5819d-8df5-4ced-b88f-9028e93d3cee
eu-west-1 - EC2Subnet - 'subnet-154d844e' - DependencyViolation: The subnet 'subnet-154d844e' has dependencies and cannot be deleted.
	status code: 400, request id: 237186aa-b035-4f64-a6e3-518bed64e240
eu-west-1 - EC2Volume - 'vol-0ddfb15461a00c3e2' - VolumeInUse: Volume vol-0ddfb15461a00c3e2 is currently attached to i-01b489457a60298dd
	status code: 400, request id: f88ff792-a17f-4fdd-9219-78a937a8d058
eu-west-1 - EC2Vpc - 'vpc-c6159fa1' - DependencyViolation: The vpc 'vpc-c6159fa1' has dependencies and cannot be deleted.
eu-west-1 - S3Object - 's3://rebuy-terraform-state-138758637120/run-terraform.lock' - triggered remove

Removal requested: 2 waiting, 6 failed, 5 skipped, 0 finished

eu-west-1 - EC2DhcpOption - 'dopt-bf2ec3d8' - DependencyViolation: The dhcpOptions 'dopt-bf2ec3d8' has dependencies and cannot be deleted.
	status code: 400, request id: d85d26e8-9f6f-42f0-811a-3b05471b0254
eu-west-1 - EC2Instance - 'i-01b489457a60298dd' - waiting
eu-west-1 - EC2KeyPair - 'test' - removed
eu-west-1 - EC2RouteTable - 'rtb-ffe91e99' - DependencyViolation: The routeTable 'rtb-ffe91e99' has dependencies and cannot be deleted.
	status code: 400, request id: adb44c0e-3f5b-4977-b2ae-7582f57fb4b7
eu-west-1 - EC2SecurityGroup - 'sg-f20f958a' - DependencyViolation: resource sg-f20f958a has a dependent object
	status code: 400, request id: c4149482-0cd2-40e0-8fa0-84a61d55a158
eu-west-1 - EC2Subnet - 'subnet-154d844e' - DependencyViolation: The subnet 'subnet-154d844e' has dependencies and cannot be deleted.
	status code: 400, request id: ba0649ba-3be8-41ee-ae0f-6b74a1f0a873
eu-west-1 - EC2Volume - 'vol-0ddfb15461a00c3e2' - VolumeInUse: Volume vol-0ddfb15461a00c3e2 is currently attached to i-01b489457a60298dd
	status code: 400, request id: 9ac3eac5-f1ef-4337-a780-228295a7ebc7
eu-west-1 - EC2Vpc - 'vpc-c6159fa1' - DependencyViolation: The vpc 'vpc-c6159fa1' has dependencies and cannot be deleted.
	status code: 400, request id: 89f870e9-1ffa-42be-9f73-76c29f088e1a

Removal requested: 1 waiting, 6 failed, 5 skipped, 1 finished

--- truncating long output ---
```

As you see *aws-nuke* now tries to delete all resources which aren't filtered,
without caring about the dependencies between them. This results in API errors
which can be ignored. They are displayed anyway, because it might be helpful
for debugging, if the error is not about dependencies.

*aws-nuke* retries deleting all resources until all specified ones are deleted
or until there are only resources with errors left.


## Install

### Use Released Binaries

The easiest way of installing it, is to download the latest
[release](https://github.com/rebuy-de/aws-nuke/releases) from GitHub.

### Compile from Source

To compile *aws-nuke* from source you need a working
[Golang](https://golang.org/doc/install) development environment. The sources
must be cloned to `$GOPATH/src/github.com/rebuy-de/aws-nuke`.

Also you need to install [Glide](https://glide.sh/),
[golint](https://github.com/golang/lint/) and [GNU
Make](https://www.gnu.org/software/make/).

Then you just need to run `make build` to compile a binary into the project
directory or `make install` go install *aws-nuke* into `$GOPATH/bin`. With
`make xc` you can cross compile *aws-nuke* for other platforms.

## Contact Channels

Feel free to create a GitHub Issue for any questions, bug reports or feature
requests.

## Contribute

You can contribute to *aws-nuke* by forking this repository, making your
changes and creating a Pull Request against our repository. If you are unsure
how to solve a problem or have other questions about a contributions, please
create a GitHub issue.

