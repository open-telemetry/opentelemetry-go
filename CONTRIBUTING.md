# Contributing to opentelemetry-go

The Go special interest group (SIG) meets regularly. See the
OpenTelemetry
[community](https://github.com/open-telemetry/community#golang-sdk)
repo for information on this and other language SIGs.

See the [public meeting
notes](https://docs.google.com/document/d/1A63zSWX0x2CyCK_LoNhmQC4rqhLpYXJzXbEPDUQ2n6w/edit#heading=h.9tngw7jdwd6b)
for a summary description of past meetings. To request edit access,
join the meeting or get in touch on
[Gitter](https://gitter.im/open-telemetry/opentelemetry-go).

## Development

There are some generated files checked into the repo. To make sure
that the generated files are up-to-date, run `make` (or `make
precommit` - the `precommit` target is the default).

The `precommit` target also fixes the formatting of the code and
checks the status of the go module files.

If after running `make precommit` the output of `git status` contains
`nothing to commit, working tree clean` then it means that everything
is up-to-date and properly formatted.

## Pull Requests

### How to Send Pull Requests

Everyone is welcome to contribute code to `opentelemetry-go` via
GitHub pull requests (PRs).

To create a new PR, fork the project in GitHub and clone the upstream
repo:

```sh
$ go get -d go.opentelemetry.io/otel
```

(This may print some warning about "build constraints exclude all Go
files", just ignore it.)

This will put the project in `${GOPATH}/src/go.opentelemetry.io/otel`. You
can alternatively use `git` directly with:

```sh
$ git clone https://github.com/open-telemetry/opentelemetry-go
```

(Note that `git clone` is *not* using the `go.opentelemetry.io/otel` name -
that name is a kind of a redirector to GitHub that `go get` can
understand, but `git` does not.)

This would put the project in the `opentelemetry-go` directory in
current working directory.

Enter the newly created directory and add your fork as a new remote:

```sh
$ git remote add <YOUR_FORK> git@github.com:<YOUR_GITHUB_USERNAME>/opentelemetry-go
```

Check out a new branch, make modifications, run linters and tests, and
push the branch to your fork:

```sh
$ git checkout -b <YOUR_BRANCH_NAME>
# edit files
$ make precommit
$ make test
$ git add -p
$ git commit
$ git push <YOUR_FORK> <YOUR_BRANCH_NAME>
```

Open a pull request against the main `opentelemetry-go` repo.

### How to Receive Comments

* If the PR is not ready for review, please put `[WIP]` in the title,
  tag it as `work-in-progress`, or mark it as
  [`draft`](https://github.blog/2019-02-14-introducing-draft-pull-requests/).
* Make sure CLA is signed and CI is clear.

### How to Get PRs Merged

A PR is considered to be **ready to merge** when:

* It has received two approvals from Collaborators/Maintainers (at
  different companies).
* Major feedbacks are resolved.
* It has been open for review for at least one working day. This gives
  people reasonable time to review.
* Trivial change (typo, cosmetic, doc, etc.) doesn't have to wait for
  one day.
* Urgent fix can take exception as long as it has been actively
  communicated.

Any Collaborator/Maintainer can merge the PR once it is **ready to
merge**.

## Design Choices

As with other OpenTelemetry clients, opentelemetry-go follows the
[opentelemetry-specification](https://github.com/open-telemetry/opentelemetry-specification).

It's especially valuable to read through the [library
guidelines](https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/library-guidelines.md).

### Focus on Capabilities, Not Structure Compliance

OpenTelemetry is an evolving specification, one where the desires and
use cases are clear, but the method to satisfy those uses cases are
not.

As such, Contributions should provide functionality and behavior that
conforms to the specification, but the interface and structure is
flexible.

It is preferable to have contributions follow the idioms of the
language rather than conform to specific API names or argument
patterns in the spec.

For a deeper discussion, see:
https://github.com/open-telemetry/opentelemetry-specification/issues/165

## Style Guide

* Make sure to run `make precommit` - this will find and fix the code
  formatting.

## Become an Approver or a Maintainer

See the [community membership document in OpenTelemetry community
repo](https://github.com/open-telemetry/community/blob/master/community-membership.md).
