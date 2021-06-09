# Contributing

- [Contributing](#contributing)
  - [Architecture Decision Records (ADR)](#architecture-decision-records-adr)
  - [Pull Requests](#pull-requests)
    - [Process for reviewing PRs](#process-for-reviewing-prs)
    - [Updating Documentation](#updating-documentation)
  - [Forking](#forking)
  - [Dependencies](#dependencies)
  - [Protobuf](#protobuf)
  - [Testing](#testing)
  - [Branching Model and Release](#branching-model-and-release)
    - [PR Targeting](#pr-targeting)
    - [Development Procedure](#development-procedure)
    - [Pull Merge Procedure](#pull-merge-procedure)
    - [Release Procedure](#release-procedure)
    - [Point Release Procedure](#point-release-procedure)
  - [Code Owner Membership](#code-owner-membership)

Thank you for considering making contributions to lfb-sdk and related
repositories!

Contributing to this repo can mean many things such as participating in
discussion or proposing code changes. To ensure a smooth workflow for all
contributors, the general procedure for contributing has been established:

1. Either [open](https://github.com/line/lfb-sdk/issues/new/choose) or
   [find](https://github.com/line/lfb-sdk/issues) an issue you'd like to help with
2. Participate in thoughtful discussion on that issue
3. If you would like to contribute:
   1. If the issue is a proposal, ensure that the proposal has been accepted
   2. Ensure that nobody else has already begun working on this issue. If they have,
      make sure to contact them to collaborate
   3. If nobody has been assigned for the issue and you would like to work on it,
      make a comment on the issue to inform the community of your intentions
      to begin work
   4. Follow standard Github best practices: fork the repo, branch from the
      HEAD of `master`, make some commits, and submit a PR to `master`
      - For core developers working within the cosmos-sdk repo, to ensure a clear
        ownership of branches, branches must be named with the convention
        `{moniker}/{issue#}-branch-name`
   5. Be sure to submit the PR in `Draft` mode submit your PR early, even if
      it's incomplete as this indicates to the community you're working on
      something and allows them to provide comments early in the development process
   6. When the code is complete it can be marked `Ready for Review`
   7. Be sure to include a relevant change log entry in the `Unreleased` section
      of `CHANGELOG.md` (see file for log format)

Note that for very small or blatantly obvious problems (such as typos) it is
not required to an open issue to submit a PR, but be aware that for more complex
problems/features, if a PR is opened before an adequate design discussion has
taken place in a github issue, that PR runs a high likelihood of being rejected.

Other notes:

- Please make sure to run `make format` before every commit - the easiest way
  to do this is have your editor run it for you upon saving a file. Additionally
  please ensure that your code is lint compliant by running `golangci-lint run`.
  A convenience git `pre-commit` hook that runs the formatters automatically
  before each commit is available in the `contrib/githooks/` directory.

## Architecture Decision Records (ADR)

When proposing an architecture decision for the SDK, please create an [ADR](./docs/architecture/README.md)
so further discussions can be made. We are following this process so all involved parties are in
agreement before any party begins coding the proposed implementation. If you would like to see some examples
of how these are written refer to the current [ADRs](./docs/architecture)

## Pull Requests

To accommodate review process we suggest that PRs are categorically broken up.
Ideally each PR addresses only a single issue. Additionally, as much as possible
code refactoring and cleanup should be submitted as a separate PRs from bugfixes/feature-additions.

### Process for reviewing PRs

All PRs require two Reviews before merge (except docs changes, or variable name-changes which only require one). When reviewing PRs please use the following review explanations:

- `LGTM` without an explicit approval means that the changes look good, but you haven't pulled down the code, run tests locally and thoroughly reviewed it.
- `Approval` through the GH UI means that you understand the code, documentation/spec is updated in the right places, you have pulled down and tested the code locally. In addition:
  - You must also think through anything which ought to be included but is not
  - You must think through whether any added code could be partially combined (DRYed) with existing code
  - You must think through any potential security issues or incentive-compatibility flaws introduced by the changes
  - Naming must be consistent with conventions and the rest of the codebase
  - Code must live in a reasonable location, considering dependency structures (e.g. not importing testing modules in production code, or including example code modules in production code).
  - if you approve of the PR, you are responsible for fixing any of the issues mentioned here and more
- If you sat down with the PR submitter and did a pairing review please note that in the `Approval`, or your PR comments.
- If you are only making "surface level" reviews, submit any notes as `Comments` without adding a review.

### Updating Documentation

If you open a PR on the LFB SDK, it is mandatory to update the relevant documentation in /docs.

- If your change relates to the core SDK (baseapp, store, ...), please update the `docs/basics/`, `docs/core/` and/or `docs/building-modules/` folders.
- If your changes relate to the core of the CLI or Light-client (not specifically to module's CLI/Rest), please modify the `docs/interfaces/` folder.
- If your changes relate to a module, please update the module's spec in `x/moduleName/docs/spec/`.

## Forking

Please note that Go requires code to live under absolute paths, which complicates forking.
While my fork lives at `https://github.com/someone/lfb-sdk`,
the code should never exist at `$GOPATH/src/github.com/someone/lfb-sdk`.
Instead, we use `git remote` to add the fork as a new remote for the original repo,
`$GOPATH/src/github.com/line/lfb-sdk`, and do all the work there.

For instance, to create a fork and work on a branch of it, I would:

- Create the fork on github, using the fork button.
- Go to the original repo checked out locally (i.e. `$GOPATH/src/github.com/line/lfb-sdk`)
- `git remote rename origin upstream`
- `git remote add origin git@github.com:someone/lfb-sdk.git`

Now `origin` refers to my fork and `upstream` refers to the lfb-sdk version.
So I can `git push -u origin master` to update my fork, and make pull requests to lfb-sdk from there.
Of course, replace `someone` with your git handle.

To pull in updates from the origin repo, run

- `git fetch upstream`
- `git rebase upstream/master` (or whatever branch you want)

Please don't make Pull Requests from `master`.

## Dependencies

We use [Go 1.15 Modules](https://github.com/golang/go/wiki/Modules) to manage
dependency versions.

The master branch of every LFB repository should just build with `go get`,
which means they should be kept up-to-date with their dependencies, so we can
get away with telling people they can just `go get` our software.

Since some dependencies are not under our control, a third party may break our
build, in which case we can fall back on `go mod tidy -v`.

## Protobuf

We use [Protocol Buffers](https://developers.google.com/protocol-buffers) along with [gogoproto](https://github.com/gogo/protobuf) to generate code for use in lfb-sdk.

For determinstic behavior around Protobuf tooling, everything is containerized using Docker. Make sure to have Docker installed on your machine, or head to [Docker's website](https://docs.docker.com/get-docker/) to install it.

For formatting code in `.proto` files, you can run `make proto-format` command.

For linting and checking breaking changes, we use [buf](https://buf.build/). You can use the commands `make proto-lint` and `make proto-check-breaking` to respectively lint your proto files and check for breaking changes.

To generate the protobuf stubs, you can run `make proto-gen`.

We also added the `make proto-all` command to run all the above commands sequentially.

In order for imports to properly compile in your IDE, you may need to manually set your protobuf path in your IDE's workspace settings/config.

For example, in vscode your `.vscode/settings.json` should look like:

```
{
    "protoc": {
        "options": [
        "--proto_path=${workspaceRoot}/proto",
        "--proto_path=${workspaceRoot}/third_party/proto"
        ]
    }
}
```

## Testing

Tests can be ran by running `make test` at the top level of the SDK repository. 

We expect tests to use `require` or `assert` rather than `t.Skip` or `t.Fail`,
unless there is a reason to do otherwise.
When testing a function under a variety of different inputs, we prefer to use
[table driven tests](https://github.com/golang/go/wiki/TableDrivenTests).
Table driven test error messages should follow the following format
`<desc>, tc #<index>, i #<index>`.
`<desc>` is an optional short description of whats failing, `tc` is the
index within the table of the testcase that is failing, and `i` is when there
is a loop, exactly which iteration of the loop failed.
The idea is you should be able to see the
error message and figure out exactly what failed.
Here is an example check:

```go
<some table>
for tcIndex, tc := range cases {
  <some code>
  for i := 0; i < tc.numTxsToTest; i++ {
      <some code>
			require.Equal(t, expectedTx[:32], calculatedTx[:32],
				"First 32 bytes of the txs differed. tc #%d, i #%d", tcIndex, i)
```

## Branching Model and Release

TBD

### PR Targeting

TBD

### Development Procedure

TBD

### Pull Merge Procedure

TBD

### Release Procedure

TBD

### Point Release Procedure

TBD

## Code Owner Membership

TBD
