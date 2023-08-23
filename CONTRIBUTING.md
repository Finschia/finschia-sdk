# Contributing

- [Contributing](#contributing)
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

Thank you for considering making contributions to finschia-sdk and related
repositories!

Contributing to this repo can mean many things such as participating in
discussion or proposing code changes. To ensure a smooth workflow for all
contributors, the general procedure for contributing has been established:

1. Either [open](https://github.com/Finschia/finschia-sdk/issues/new/choose) or
   [find](https://github.com/Finschia/finschia-sdk/issues) an issue you'd like to help with
2. Participate in thoughtful discussion on that issue
3. If you would like to contribute:
   1. If the issue is a proposal, ensure that the proposal has been accepted
   2. Ensure that nobody else has already begun working on this issue. If they have,
      make sure to contact them to collaborate
   3. If nobody has been assigned for the issue and you would like to work on it,
      make a comment on the issue to inform the community of your intentions
      to begin work
   4. Follow standard Github best practices: fork the repo, branch from the
      HEAD of `main`, make some commits, and submit a PR to `main`
      - For core developers working within the finschia-sdk repo, to ensure a clear
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

## Commit Convention
All commits pushed to `finschia-sdk` are classified according to the rules below.
Therefore, it cannot be included in the release notes unless you properly assign the prefix to the commit. For example, you can write commit message like `feat: this is new feature`. However, if you want to clarify the scope of the commit, you can write it as follows `feat(x/token): this is new feature`.
- `chore` : a commit of the type chore makes trivial changes. If you don't want to include specific commit in the release note, add this prefix.
- `feat` : a commit of the type feat introduces a new feature to the codebase
- `fix` : a commit of the type fix patches a bug in your codebase
- `refactor` : A code change that neither fixes a bug nor adds a feature
- `BREAKING CHANGE` : a commit of the type breaking changes to the codebase. This changes breaks down the backward compatibility
- `perf` : A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `docs`: Documentation only changes
- `build`: Changes that affect the build system or external dependencies
- `ci`: Changes to our CI configuration files and scripts

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

If you open a PR on the LBM SDK, it is mandatory to update the relevant documentation in /docs.

- If your change relates to the core SDK (baseapp, store, ...), please update the `docs/basics/`, `docs/core/` and/or `docs/building-modules/` folders.
- If your changes relate to the core of the CLI (not specifically to module's CLI/Rest), please modify the `docs/run-node/` folder.
- If your changes relate to a module, please update the module's spec in `x/moduleName/docs/spec/`.

## Forking

Please note that Go requires code to live under absolute paths, which complicates forking.
While my fork lives at `https://github.com/someone/finschia-sdk`,
the code should never exist at `$GOPATH/src/github.com/someone/finschia-sdk`.
Instead, we use `git remote` to add the fork as a new remote for the original repo,
`$GOPATH/src/github.com/Finschia/finschia-sdk`, and do all the work there.

For instance, to create a fork and work on a branch of it, I would:

- Create the fork on github, using the fork button.
- Go to the original repo checked out locally (i.e. `$GOPATH/src/github.com/Finschia/finschia-sdk`)
- `git remote rename origin upstream`
- `git remote add origin git@github.com:someone/finschia-sdk.git`

Now `origin` refers to my fork and `upstream` refers to the finschia-sdk version.
So I can `git push -u origin main` to update my fork, and make pull requests to finschia-sdk from there.
Of course, replace `someone` with your git handle.

To pull in updates from the origin repo, run

- `git fetch upstream`
- `git rebase upstream/main` (or whatever branch you want)

Please don't make Pull Requests from `main`.

## Dependencies

We use [Go 1.20 Modules](https://github.com/golang/go/wiki/Modules) to manage
dependency versions.

The `main` branch of every LBM repository should just build with `go get`,
which means they should be kept up-to-date with their dependencies, so we can
get away with telling people they can just `go get` our software.

Since some dependencies are not under our control, a third party may break our
build, in which case we can fall back on `go mod tidy -v`.

## Protobuf

We use [Protocol Buffers](https://developers.google.com/protocol-buffers) along with [gogoproto](https://github.com/gogo/protobuf) to generate code for use in finschia-sdk.

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

User-facing repos should adhere to the trunk based development branching model: https://trunkbaseddevelopment.com/.

Libraries need not follow the model strictly, but would be wise to.

The SDK utilizes [semantic versioning](https://semver.org/).

### PR Targeting

Ensure that you base and target your PR on the `main` branch.

All feature additions should be targeted against `main`. Bug fixes for an outstanding release candidate
should be targeted against the release candidate branch.

### Development Procedure

- the latest state of development is on `main`
- `main` must never fail `make lint test test-race`
- `main` should not fail `make lint`
- no `--force` onto `main` (except when reverting a broken commit, which should seldom happen)
- create a development branch either on github.com/Finschia/finschia-sdk, or your fork (using `git remote add origin`)
- before submitting a pull request, begin `git rebase` on top of `main`

### Pull Merge Procedure

- ensure pull branch is rebased on `main`
- run `make test` to ensure that all tests pass
- merge pull request (We are using `squash and merge` for small features)

### Release Procedure

- Start on `main`
- Create the release candidate branch `rc/v*` (going forward known as **RC**)
  and ensure it's protected against pushing from anyone except the release
  manager/coordinator
  - **no PRs targeting this branch should be merged unless exceptional circumstances arise**
- On the `RC` branch, prepare a new version section in the `CHANGELOG.md`
  - All links must be link-ified: `$ python ./scripts/linkify_changelog.py CHANGELOG.md`
  - Copy the entries into a `RELEASE_CHANGELOG.md`, this is needed so the bot knows which entries to add to the release page on github.
- Kick off a large round of simulation testing (e.g. 400 seeds for 2k blocks)
- If errors are found during the simulation testing, commit the fixes to `main`
  and create a new `RC` branch (making sure to increment the `rcN`)
- After simulation has successfully completed, create the release branch
  (`release/vX.XX.X`) from the `RC` branch
- Create a PR to `main` to incorporate the `CHANGELOG.md` updates
- Tag the release (use `git tag -a`) and create a release in Github
- Delete the `RC` branches

