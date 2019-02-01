# Contributing

Thank you for wanting to contribute to *aws-nuke*.

Because of the amount of AWS services and their rate of change, we rely on your
participation. For the same reason we can only act retroactive on changes of
AWS services. Otherwise it would be a fulltime job to keep up with AWS.


## How Can I Contribute?

### Some Resource Is Not Supported by *aws-nuke*

If a resource is not yet supported by *aws-nuke*, you have two options to
resolve this:

* File [an issue](https://github.com/rebuy-de/aws-nuke/issues/new) and describe
  which resource is missing. This way someone can take care of it.
* Add the resource yourself and open a Pull Request. Please follow the
  guidelines below to see how to create such a resource.


### Some Resource Does Not Get Deleted

Please check the following points before creating a bug issue:

* Is the resource actually supported by *aws-nuke*? If not, please follow the
  guidelines above.
* Are there permission problems? In this case *aws-nuke* will print errors
  that usually contain the status code `403`.
* Did you just get scared by an error that was printed? *aws-nuke* does not
  know about dependencies between resources. To work around this it will just
  retry deleting all resources in multiple iterations. Therefore it is normal
  that there are a lot of dependency errors in the first one. The iterations
  are separated by lines starting with `Removal requested: ` and only the
  errors in the last block indicate actual errros.

File [an issue](https://github.com/rebuy-de/aws-nuke/issues/new) and describe
as accurately as possible how to generate the resource on AWS that cause the
errors in *aws-nuke*. Ideally this is provided in a reproducible way like
a Terraform template or AWS CLI commands.


### I Have Ideas to Improve *aws-nuke*

You should take these steps if you have an idea how to improve *aws-nuke*:

1. Check the [issues page](https://github.com/rebuy-de/aws-nuke/issues),
   whether someone already had the same or a similar idea.
2. Also check the [closed
   issues](https://github.com/rebuy-de/aws-nuke/issues?utf8=%E2%9C%93&q=is%3Aissue),
   because this might have already been implemented, but not yet released. Also
   the idea might not be viable for unobvious reasons.
3. Join the discussion, if there is already an related issue. If this is not
   the case, open a new issue and describe your idea. Afterwards, we can
   discuss this idea and form a proposal.


### I Just Have a Question

Please use our mailing list for questions: aws-nuke@googlegroups.com. You can
also search in the mailing list archive, whether someone already had the same
problem: https://groups.google.com/d/forum/aws-nuke


## Resource Guidelines

### Consider Pagination

Most AWS resources are paginated and all resources should handle that.


### Use Properties Instead of String Functions

Currently, each resource can offer two functions to describe itself, that are
used by the user to identify it and by *aws-nuke* to filter it.

The String function is deprecated:

```go
String() string
```

The Properties function should be used instead:

```go
Properties() types.Properties
```

The interface for the String function is still there, because not all resources
are migrated yet. Please use the Properties function for new resources.


### Filter Resources That Cannot Get Removed

Some AWS APIs list resources, that cannot be deleted. For example:

* Resources that are already deleted, but still listed for some time (eg EC2 Instances).
* Resources that are created by AWS, but cannot be deleted by the user (eg some IAM Roles).

Those resources should be excluded in the filter step, rather than in the list step.


## Styleguide

### Go

#### Code Format

Like almost all Go projects, we are using `go fmt` as a single source of truth
for formatting the source code. Please use `go fmt` before committing any
change.


### Git

#### Squash

We want to keep the git history of *aws-nuke* clean. Generally, it is fine to
have one commit per resource, putting them into a single commit is also
fine. But having a lot of meaningless "fixup" and "change" commits would
clutter up the history. Therefore those commits should be squashed into a single
commit.

To squash all commits in your branch, execute:

```
git rebase --fork-point -i master
```

This opens an editor, where you should replace `pick` with `squash` on every
but the first commit. Then *git* opens another editor, where you have to update
the commit message. This message should be updated properly.

Afterwards you have to force push your changes.

> **Note:** It is generally not advised to do a history rewrite, but we
> consider that a branch and PR is owned by the author, until it gets merged.
> Therefore the author can always rewrite its own branch. An important
> implication is, that you should never add a commit to a branch of another
> author, without communicating this beforehand.


#### Rebase

We want to keep the git history of *aws-nuke* clean. Using a merge from the
master branch to update a feature branch would add unnecessary commits to the
history. Therefore, please use rebase, rather than merge to update your branch.

> **Note:** We cannot use the GitHub Squash Merge, since it would alter the
> commit author. We do not want this, because it would not properly acknowledge
> the contributions of the community.

It is recommended to [squash](#squash) your branch before doing a rebase, so you avoid
unnecessary conflicts.

To rebase your branch, simply update your master branch and run this command:

```
git rebase master
```

Afterwards you have to force push your changes.


#### Repair a Broken Branch

Sometimes, wrong *git* commands break a branch in a way that makes it hard to
properly clean it up. To fix this, we can create a new branch and put all
changes of the broken branch there.

As a first step you have to commit all changes of your broken branch. The
commit message does not matter.

Afterwards you need to create a new branch, based on `master`:

```
git checkout master
git checkout -b repair-branch
```

Then you need to merge the changes of the broken branch into the new one,
without taking over the commits itself:

```
git merge --squash broken-branch
```

Since this actually does not create a commit, you have to commit the changes
manually and write a proper commit message:

```
git commit
```

To verify that there are no unwanted changes you can do a diff to the broken
branch. This should not print any changes.

```
git diff broken-branch
```

Additionally you should test that your branch is working as expected.

If you are sure, that you new branch is working and contains all changes you
did, you can rewrite the broken branch. This avoids having to create a new Pull
Request. Be aware, that this will overwrite the broken branch.

```
# git checkout -B broken-branch
```

Afterwards you have to force push your changes.

> **Note:** If you accidentally overwrote your branch, you might be able to
> recover them with `git reflog`.


#### Setup Email

We prefer having the commit linked to the GitHub account, that is creating the
Pull Request. To make this happen, *git* must be configured with an email, that
is registered with a GitHub account.

To set the email for all git commits, you can use this command:

```
git config --global user.email "email@example.com"
```

If you want to change the email only for the *aws-nuke* repository, you can
skip the `--global` flag. You have to make sure that you are executing this in
the *aws-nuke* directory:

```
git config user.email "email@example.com"
```

If you already committed something with a wrong email, you can use this command:

```
git commit --amend --author="Author Name <email@address.com>"
```

This changes the email of the lastest commit. If you have multiple commits in
your branch, please follow the [Squash guidelines](#squash) and change the
author afterwards.
