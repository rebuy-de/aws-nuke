# aws-nuke

[![Build Status](https://travis-ci.org/rebuy-de/aws-nuke.svg?branch=master)](https://travis-ci.org/rebuy-de/aws-nuke)
[![license](https://img.shields.io/github/license/rebuy-de/aws-nuke.svg)](https://github.com/rebuy-de/aws-nuke/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/rebuy-de/aws-nuke.svg)](https://github.com/rebuy-de/aws-nuke/releases)

Remove all resources from an AWS account.

> **Development Status** *aws-nuke* is stable, but it is likely that not all AWS
resources are covered by it. Be encouraged to add missing resources and create
a Pull Request or to create an [Issue](https://github.com/rebuy-de/aws-nuke/issues/new).

## Caution!

Be aware that *aws-nuke* is a very destructive tool, hence you have to be very
careful while using it. Otherwise you might delete production data.

**We strongly advice you to not run this application on any AWS account, where
you cannot afford to lose all resources.**

To reduce the blast radius of accidents, there are some safety precautions:

1. By default *aws-nuke* only lists all nukeable resources. You need to add
   `--no-dry-run` to actually delete resources.
2. *aws-nuke* asks you twice to confirm the deletion by entering the account
   alias. The first time is directly after the start and the second time after
   listing all nukeable resources.
3. To avoid just displaying a account ID, which might gladly be ignored by
   humans, it is required to actually set an [Account
   Alias](https://docs.aws.amazon.com/IAM/latest/UserGuide/console_account-alias.html)
   for your account. Otherwise *aws-nuke* will abort.
4. The Account Alias must not contain the string `prod`. This string is
   hardcoded and it is recommended to add it to every actual production account
   (eg `mycompany-production-ecr`).
5. The config file contains a blacklist field. If the Account ID of the account
   you want to nuke is part of this blacklist, *aws-nuke* will abort. It is
   recommended, that you add every production account to this blacklist.
6. To ensure you don't just ignore the blacklisting feature, the blacklist must
   contain at least one Account ID.
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
- global

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

eu-west-1 - EC2DHCPOption - 'dopt-bf2ec3d8' - would remove
eu-west-1 - EC2Instance - 'i-01b489457a60298dd' - would remove
eu-west-1 - EC2KeyPair - 'test' - would remove
eu-west-1 - EC2NetworkACL - 'acl-6482a303' - cannot delete default VPC
eu-west-1 - EC2RouteTable - 'rtb-ffe91e99' - would remove
eu-west-1 - EC2SecurityGroup - 'sg-220e945a' - cannot delete group 'default'
eu-west-1 - EC2SecurityGroup - 'sg-f20f958a' - would remove
eu-west-1 - EC2Subnet - 'subnet-154d844e' - would remove
eu-west-1 - EC2Volume - 'vol-0ddfb15461a00c3e2' - would remove
eu-west-1 - EC2VPC - 'vpc-c6159fa1' - would remove
eu-west-1 - IAMUserAccessKey - 'my-user -> ABCDEFGHIJKLMNOPQRST' - would remove
eu-west-1 - IAMUserPolicyAttachment - 'my-user -> AdministratorAccess' - [UserName: "my-user", PolicyArn: "arn:aws:iam::aws:policy/AdministratorAccess", PolicyName: "AdministratorAccess"] - would remove
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

eu-west-1 - EC2DHCPOption - 'dopt-bf2ec3d8' - would remove
eu-west-1 - EC2Instance - 'i-01b489457a60298dd' - would remove
eu-west-1 - EC2KeyPair - 'test' - would remove
eu-west-1 - EC2NetworkACL - 'acl-6482a303' - cannot delete default VPC
eu-west-1 - EC2RouteTable - 'rtb-ffe91e99' - would remove
eu-west-1 - EC2SecurityGroup - 'sg-220e945a' - cannot delete group 'default'
eu-west-1 - EC2SecurityGroup - 'sg-f20f958a' - would remove
eu-west-1 - EC2Subnet - 'subnet-154d844e' - would remove
eu-west-1 - EC2Volume - 'vol-0ddfb15461a00c3e2' - would remove
eu-west-1 - EC2VPC - 'vpc-c6159fa1' - would remove
eu-west-1 - IAMUserAccessKey - 'my-user -> ABCDEFGHIJKLMNOPQRST' - filtered by config
eu-west-1 - IAMUserPolicyAttachment - 'my-user -> AdministratorAccess' - [UserName: "my-user", PolicyArn: "arn:aws:iam::aws:policy/AdministratorAccess", PolicyName: "AdministratorAccess"] - would remove
eu-west-1 - IAMUser - 'my-user' - filtered by config
Scan complete: 13 total, 8 nukeable, 5 filtered.

Do you really want to nuke these resources on the account with the ID 000000000000 and the alias 'aws-nuke-example'?
Do you want to continue? Enter account alias to continue.
> aws-nuke-example

eu-west-1 - EC2DHCPOption - 'dopt-bf2ec3d8' - failed
eu-west-1 - EC2Instance - 'i-01b489457a60298dd' - triggered remove
eu-west-1 - EC2KeyPair - 'test' - triggered remove
eu-west-1 - EC2RouteTable - 'rtb-ffe91e99' - failed
eu-west-1 - EC2SecurityGroup - 'sg-f20f958a' - failed
eu-west-1 - EC2Subnet - 'subnet-154d844e' - failed
eu-west-1 - EC2Volume - 'vol-0ddfb15461a00c3e2' - failed
eu-west-1 - EC2VPC - 'vpc-c6159fa1' - failed
eu-west-1 - S3Object - 's3://rebuy-terraform-state-138758637120/run-terraform.lock' - triggered remove

Removal requested: 2 waiting, 6 failed, 5 skipped, 0 finished

eu-west-1 - EC2DHCPOption - 'dopt-bf2ec3d8' - failed
eu-west-1 - EC2Instance - 'i-01b489457a60298dd' - waiting
eu-west-1 - EC2KeyPair - 'test' - removed
eu-west-1 - EC2RouteTable - 'rtb-ffe91e99' - failed
eu-west-1 - EC2SecurityGroup - 'sg-f20f958a' - failed
eu-west-1 - EC2Subnet - 'subnet-154d844e' - failed
eu-west-1 - EC2Volume - 'vol-0ddfb15461a00c3e2' - failed
eu-west-1 - EC2VPC - 'vpc-c6159fa1' - failed

Removal requested: 1 waiting, 6 failed, 5 skipped, 1 finished

--- truncating long output ---
```

As you see *aws-nuke* now tries to delete all resources which aren't filtered,
without caring about the dependencies between them. This results in API errors
which can be ignored. These errors are shown at the end of the *aws-nuke* run,
if they keep to appear.

*aws-nuke* retries deleting all resources until all specified ones are deleted
or until there are only resources with errors left.

### AWS Credentials

There are two ways to authenticate *aws-nuke*. There are static credentials and
profiles. The later one can be configured in the shared credentials file (ie
`~/.aws/credentials`) or the shared config file (ie `~/.aws/config`).

To use *static credentials* the command line flags `--access-key-id` and
`--secret-access-key` are required. The flag `--session-token` is only required
for temporary sessions.

To use *shared profiles* the command line flag `--profile` is required. The
profile must be either defined with static credentials in the [shared
credential
file](https://docs.aws.amazon.com/cli/latest/userguide/cli-multiple-profiles.html)
or in [shared config
file](https://docs.aws.amazon.com/cli/latest/userguide/cli-roles.html) with an
assuming role.

### Specifying Resource Types to Delete

*aws-nuke* deletes a lot of resources and there might be added more at any
release. Eventually, every resources should get deleted. You might want to
restrict which resources to delete. There are multiple ways to configure this.

One way are filters, which already got mentioned. This requires to know the
identifier of each resource. It is also possible to prevent whole resource
types (eg `S3Bucket`) from getting deleted with two methods.

* The `--target` flag limits nuking to the specified resource types.
* The `--exclude` flag prevent nuking of the specified resource types.

It is also possible to configure the resource types in the config file like in
these examples:

```
---
regions:
  - "eu-west-1"
account-blacklist:
- 1234567890

resource-types:
  # only nuke these three resources
  targets:
  - S3Object
  - S3Bucket
  - IAMRole

accounts:
  555133742: {}
```

```
---
regions:
  - "eu-west-1"
account-blacklist:
- 1234567890

resource-types:
  # don't nuke IAM users
  excludes:
  - IAMUser

accounts:
  555133742: {}
```

If targets are specified in multiple places (eg CLI and account specific), then
a resource type must be specified in all places. In other words each
configuration limits the previous ones.

If an exclude is used, then all its resource types will not be deleted.

**Hint:** You can see all available resource types with this command:

```
aws-nuke resource-types
```


### Feature Flags

There are some features, which are quite opinionated. To make those work for
everyone, *aws-nuke* has flags to manually enable those features. These can be
configured on the root-level of the config, like this:

```yaml
---
feature-flags:
  disable-deletion-protection:
    RDSInstance: true
```


### Filtering Resources

It is possible to filter this is important for not deleting the current user
for example or for resources like S3 Buckets which have a globally shared
namespace and might be hard to recreate. Currently the filtering is based on
the resource identifier. The identifier will be printed as the first step of
*aws-nuke* (eg `i-01b489457a60298dd` for an EC2 instance).

**Note: Even with filters you should not run aws-nuke on any AWS account, where
you cannot afford to lose all resources. It is easy to make mistakes in the
filter configuration. Also, since aws-nuke is in continous development, there
is always a possibility to introduce new bugs, no matter how careful we review
new code.**

The filters are part of the account-specific configuration and are grouped by
resource types. This is an example of a config that deletes all resources but
the `admin` user with its access permissions and two access keys:

```yaml
---
regions:
- global
- eu-west-1

account-blacklist:
- 1234567890

accounts:
  0987654321:
    filters:
      IAMUser:
      - "admin"
      IAMUserPolicyAttachment:
      - "admin -> AdministratorAccess"
      IAMUserAccessKey:
      - "admin -> AKSDAFRETERSDF"
      - "admin -> AFGDSGRTEWSFEY"
```

Any resource whose resource identifier exactly matches any of the filters in
the list will be skipped. These will be marked as "filtered by config" on the
*aws-nuke* run.

#### Filter Properties

Some resources support filtering via properties. When a resource support these
properties, they will be listed in the output like in this example:

```
global - IAMUserPolicyAttachment - 'admin -> AdministratorAccess' - [RoleName: "admin", PolicyArn: "arn:aws:iam::aws:policy/AdministratorAccess", PolicyName: "AdministratorAccess"] - would remove
```

To use properties, it is required to specify a object with `properties` and
`value` instead of the plain string.

These types can be used to simplify the configuration. For example, it is
possible to protect all access keys of a single user:

```yaml
IAMUserAccessKey:
- property: UserName
  value: "admin"
```

#### Filter Types

There are also additional comparision types than an exact match:

* `exact` – The identifier must exactly match the given string. This is the default.
* `contains` – The identifier must contain the given string.
* `glob` – The identifier must match against the given [glob
  pattern](https://en.wikipedia.org/wiki/Glob_(programming)). This means the
  string might contains wildcards like `*` and `?`. Note that globbing is
  designed for file paths, so the wildcards do not match the directory
  separator (`/`). Details about the glob pattern can be found in the [library
  documentation](https://godoc.org/github.com/mb0/glob).
* `regex` – The identifier must match against the given regular expression.
  Details about the syntax can be found in the [library
  documentation](https://golang.org/pkg/regexp/syntax/).

To use a non-default comparision type, it is required to specify a object with
`type` and `value` instead of the plain string.

These types can be used to simplify the configuration. For example, it is
possible to protect all access keys of a single user by using `glob`:

```yaml
IAMUserAccessKey:
- type: glob
  value: "admin -> *"
```


#### Using Them Together

It is also possible to use Filter Properties and Filter Types together. For
example to protect all Hosted Zone of a specific TLD:

```yaml
Route53HostedZone:
- property: Name
  type: glob
  value: "*.rebuy.cloud."
```

####  Inverting Filter Results

Any filter result can be inverted by using `invert: true`, for example:
```yaml
CloudFormationStack:
- property: Name
  value: "foo"
  invert: true
```

In this case *any* CloudFormationStack ***but*** the ones called "foo" will be
filtered. Be aware that *aws-nuke* internally takes every resource and applies
every filter on it. If a filter matches, it marks the node as filtered.


#### Filter Presets

It might be the case that some filters are the same across multiple accounts.
This especially could happen, if provisioning tools like Terraform are used or
if IAM resources follow the same pattern.

For this case *aws-nuke* supports presets of filters, that can applied on
multiple accounts. A configuration could look like this:

```yaml
---
regions:
- "global"
- "eu-west-1"

account-blacklist:
- 1234567890

accounts:
  555421337:
    presets:
    - "common"
  555133742:
    presets:
    - "common"
    - "terraform"
  555134237:
    presets:
    - "common"
    - "terraform"
    filters:
      EC2KeyPair:
      - "notebook"

presets:
  terraform:
    filters:
      S3Bucket:
      - type: glob
        value: "my-statebucket-*"
      DynamoDBTable:
      - "terraform-lock"
  common:
    filter:
      IAMRole:
      - "OrganizationAccountAccessRole"
```


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

### Docker

You can run *aws-nuke* with Docker by using a command like this:

```bash
$ docker run \
    --rm -it \
    -v /full-path/to/nuke-config.yml:/home/aws-nuke/config.yml \
    -v /home/user/.aws:/home/aws-nuke/.aws \
    quay.io/rebuy/aws-nuke:v2.11.0 \
    --profile default \
    --config /home/aws-nuke/config.yml
```

To make it work, you need to adjust the paths for the AWS config and the
*aws-nuke* config.

Also you need to specify the correct AWS profile. Instead of mounting the AWS
directory, you can use the `--access-key-id` and `--secret-access-key` flags.

Make sure you use the latest version in the image tag. Alternatiely you can use
`master` for the latest development version, but be aware that this is more
likely to break at any time.


## Contact Channels

Feel free to create a GitHub Issue for any questions, bug reports or feature
requests.

## Contribute

You can contribute to *aws-nuke* by forking this repository, making your
changes and creating a Pull Request against our repository. If you are unsure
how to solve a problem or have other questions about a contributions, please
create a GitHub issue.
