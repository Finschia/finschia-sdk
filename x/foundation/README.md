<!--
order: 0
title: Foundation Overview
parent:
  title: "foundation"
-->

# `x/foundation`

## Abstract

This module provides the functionalities related to the foundation.
The foundation can turn off these functionalities irreversibly, through the
corresponding proposal. Therefore, the users can ensure that no one can bring
back these foundation-specific functionalities.

## Contents

* [Concepts](#concepts)
* [State](#state)
* [Msg Service](#msg-service)
    * [Msg/UpdateParams](#msgupdateparams)
    * [Msg/UpdateDecisionPolicy](#msgupdatedecisionpolicy)
    * [Msg/UpdateMembers](#msgupdatemembers)
    * [Msg/LeaveFoundation](#msgleavefoundation)
    * [Msg/SubmitProposal](#msgsubmitproposal)
    * [Msg/WithdrawProposal](#msgwithdrawproposal)
    * [Msg/Vote](#msgvote)
    * [Msg/Exec](#msgexec)
    * [Msg/Grant](#msggrant)
    * [Msg/Revoke](#msgrevoke)
    * [Msg/FundTreasury](#msgfundtreasury)
    * [Msg/WithdrawFromTreasury](#msgwithdrawfromtreasury)
    * [Msg/GovMint](#msggovmint)
* [Events](#events)
    * [EventUpdateFoundationParams](#eventupdatefoundationparams)
    * [EventUpdateDecisionPolicy](#eventupdatedecisionpolicy)
    * [EventUpdateMembers](#eventupdatedmembers)
    * [EventLeaveFoundation](#eventleavefoundation)
    * [EventSubmitProposal](#eventsubmitproposal)
    * [EventWithdrawProposal](#eventwithdrawproposal)
    * [EventVote](#eventvote)
    * [EventExec](#eventexec)
    * [EventGrant](#eventgrant)
    * [EventRevoke](#eventrevoke)
    * [EventFundTreasury](#eventfundedtreasury)
    * [EventWithdrawFromTreasury](#eventwithdrawedfromtreasury)
    * [EventGovMint](#eventgovmint)
* [Client](#client)
    * [CLI](#cli)
    * [gRPC](#grpc)

# Concepts

## Authority

`x/foundation`'s authority is a module account associated with the foundation
and a decision policy. It is an "administrator" which has the ability to add,
remove and update members in the foundation.
`x/foundation` has several messages which cannot be triggered but by the
authority. It includes membership management messages, and other messages which
controls the assets of the foundation.

**Note:** The authority is a module account, which means no one has the private
key of the authority. Hence, foundation members MUST propose, vote and execute
the corresponding proposal.

## Decision Policy

A decision policy is the rules that dictate whether a proposal should pass or
not based on its tally outcome.

All decision policies generally would have a mininum execution period and a
maximum voting window. The minimum execution period is the minimum amount of
time that must pass after submission in order for a proposal to potentially be
executed, and it may be set to 0. The maximum voting window is the maximum time
after submission that a proposal may be voted on before it is tallied.

The chain developer also defines an app-wide maximum execution period, which is
the maximum amount of time after a proposal's voting period end where the
members are allowed to execute a proposal.

The current foundation module comes shipped with two decision policies:
threshold and percentage. Any chain developer can extend upon these two, by
creating custom decision policies, as long as they adhere to the
`DecisionPolicy` interface:

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/x/foundation/foundation.go#L90-L103

### Threshold decision policy

A threshold decision policy defines a threshold of yes votes (based on a tally
of voter weights) that must be achieved in order for a proposal to pass. For
this decision policy, abstain and veto are simply treated as no's.

### Percentage decision policy

A percentage decision policy is similar to a threshold decision policy, except
that the threshold is not defined as a constant weight, but as a percentage.
It's more suited for a foundation where the membership can be updated, as the
percentage threshold stays the same, and doesn't depend on how the number of
members get updated.

### Outsourcing decision policy

A outsourcing decision policy is a policy set after `x/foundation` decides to
outsource its proposal relevant features to other modules (e.g. `x/group`).
It means one can expect that any states relevant to the feature must be removed
in the update to this policy.

## Proposal

Any foundation member(s) can submit a proposal for the foundation policy
account to decide upon. A proposal consists of a set of messages that will be
executed if the proposal passes as well as any metadata associated with the
proposal.

### Voting

There are four choices to choose while voting - yes, no, abstain and veto. Not
all decision policies will take the four choices into account. Votes can
contain some optional metadata.

In the current implementation, the voting window begins as soon as a proposal
is submitted, and the end is defined by the decision policy.

### Withdrawing Proposals

Proposals can be withdrawn any time before the voting period end, either by the
module's authority or by one of the proposers. Once withdrawn, it is marked as
`PROPOSAL_STATUS_WITHDRAWN`, and no more voting or execution is allowed on it.

### Aborted Proposals

If the decision policy is updated during the voting period of the proposal,
then the proposal is marked as `PROPOSAL_STATUS_ABORTED`, and no more voting or
execution is allowed on it. This is because the decision policy defines the
rules of proposal voting and execution, so if those rules change during the
lifecycle of a proposal, then the proposal should be marked as stale.

### Tallying

Tallying is the counting of all votes on a proposal. It can be triggered by the
following two factors:

* either someone tries to execute the proposal (see next section), which can
  happen on a `Msg/Exec` transaction, or a `Msg/{SubmitProposal,Vote}`
  transaction with the `Exec` field set. When a proposal execution is
  attempted, a tally is done first to make sure the proposal passes.
* or on `EndBlock` when the proposal's voting period end just passed.

If the tally result passes the decision policy's rules, then the proposal is
marked as `PROPOSAL_STATUS_ACCEPTED`, or else it is marked as
`PROPOSAL_STATUS_REJECTED`. In any case, no more voting is allowed anymore, and
the tally result is persisted to state in the proposal's `FinalTallyResult`.

### Executing Proposals

Proposals are executed only when the tallying is done, and the decision policy
allows the proposal to pass based on the tally outcome. They are marked by the
status `PROPOSAL_STATUS_ACCEPTED`. Execution must happen before a duration of
`MaxExecutionPeriod` (set by the chain developer) after each proposal's voting
period end.

Proposals will not be automatically executed by the chain in this current
design, but rather a member must submit a `Msg/Exec` transaction to attempt to
execute the proposal based on the current votes and decision policy. Any member
can execute proposals that have been accepted, and execution fees are paid by
the proposal executor.

It's also possible to try to execute a proposal immediately on creation or on
new votes using the `Exec` field of `Msg/SubmitProposal` and `Msg/Vote`
requests. In the former case, proposers signatures are considered as yes votes.
In these cases, if the proposal can't be executed (i.e. it didn't pass the
decision policy's rules), it will still be opened for new votes and
could be tallied and executed later on.

A successful proposal execution will have its `ExecutorResult` marked as
`PROPOSAL_EXECUTOR_RESULT_SUCCESS`. The proposal will be automatically pruned
after execution. On the other hand, a failed proposal execution will be marked
as `PROPOSAL_EXECUTOR_RESULT_FAILURE`. Such a proposal can be re-executed
multiple times, until it expires after `MaxExecutionPeriod` after voting period
end.

## Pruning

Proposals and votes are automatically pruned to avoid state bloat.

Votes are pruned:

* either after a successful tally, i.e. a tally whose result passes the
  decision policy's rules, which can be trigged by a `Msg/Exec` or a
  `Msg/{SubmitProposal,Vote}` with the `Exec` field set,
* or on `EndBlock` right after the proposal's voting period end. This applies
  to proposals with status `aborted` or `withdrawn` too.
* or after updating the membership or decision policy.

whichever happens first.

Proposals are pruned:

* on `EndBlock` whose proposal status is `withdrawn` or `aborted` on proposal's
  voting period end before tallying,
* and either after a successful proposal execution,
* or on `EndBlock` right after the proposal's `voting_period_end` +
  `max_execution_period` (defined as an app-wide configuration) is passed,

whichever happens first.

## Authorization

The foundation module defines interfaces of authorizations on messages to
enforce_censorship_ on its execution. The other modules may deny the execution
of the message based on the information in the foundation.
`Authorization` is an interface that must be implemented by a concrete
authorization logic to validate and execute grants. `Authorization`s are
extensible and can be defined for any Msg service method even outside of the
module where the Msg method is defined.

**Note:** The foundation module's `Authorization` is different from that of
`x/authz`, while the latter allows an account to perform actions on behalf of
another account.

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/x/foundation/authz.go#L10-L27

## Built-in Authorizations

### ReceiveFromTreasuryAuthorization

`ReceiveFromTreasuryAuthorization` implements the `Authorization` interface for
the [Msg/WithdrawFromTreasury](#msgwithdrawfromtreasury).

**Note:** The subject which executes
`lbm.foundation.v1.MsgWithdrawFromTreasury` is the foundation.

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/proto/lbm/foundation/v1/authz.proto#L9-L13

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/x/foundation/authz.pb.go#L27-L30

### CreateValidatorAuthorization

`CreateValidatorAuthorization` implements the `Authorization` interface for the
[Msg/CreateValidator](../stakingplus/spec/03_messages.md#msgcreatevalidator).
An account must have this authorization prior to sending the message.

**Note:** You MUST provide the `CreateValidatorAuthorization`s into the genesis
if `Msg/CreateValidator` is being censored (`CensoredMsgTypeUrls` contains the
url of `Msg/CreateValidator`), or the chain cannot be started.

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/proto/lbm/stakingplus/v1/authz.proto#L9-L15

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/x/stakingplus/authz.pb.go#L27-L31

## Foundation Treasury

`x/foundation` intercepts the rewards prior to its distribution
(by `x/distribution`). The rate would be `FoundationTax`.

The foundation can withdraw coins from the treasury. The recipient must have
the corresponding authorization (`ReceiveFromTreasuryAuthorization`) prior to
sending the message `Msg/WithdrawFromTreasury`.

**Note:** After setting the tax rate to zero, you cannot set it to a non-zero
value again (irreversible), which means you must set it to a non-zero value in
the genesis to make it work.

## GovMint

When the chain is first started, it may be necessary to mint a large amount of 
coins at most once for initial validators or for specific purposes. Newly minted
coins are transferred to the treasury pool.

# State

## Params

* Params: `0x00 -> PropocolBuffer(Params)`.

### FoundationTax

The value of `FoundationTax` is the foundation tax rate.

### CensoredMsgTypeUrls

The `CensoredMsgTypeUrls` contains the urls of the messages under the
censorship.

## FoundationInfo

`FoundationInfo` contains the information relevant to the foundation.

* FoundationInfo: `0x01 -> ProtocolBuffer(FoundationInfo)`.

### Version

The `Version` is used to track changes to the foundation membership. Whenever
the membership is changed, this value is incremented, which will cause
proposals based on older versions to fail.

### TotalWeight

The `TotalWeight` is the number of the foundation members.

### DecisionPolicy

The `DecisionPolicy` is the decision policy of the foundation.

## Member

The `Member` is the foundation member.

* Member: `0x10 | []byte(member.Address) -> ProtocolBuffer(Member)`.

## PreviousProposalID

The value of the `PreviousProposalID` is the last used proposal ID. The chain
uses this value to issue the ID of the next new proposal.

* PreviousProposalID: `0x11 -> BigEndian(ProposalId)`.

## Proposal

* Proposal: `0x12 | BigEndian(ProposalId) -> ProtocolBuffer(Proposal)`.

## ProposalByVotingPeriodEnd

`ProposalByVotingPeriodEnd` allows to retrieve proposals sorted by
chronological `voting_period_end`. This index is used when tallying the
proposal votes at the end of the voting period, and for pruning proposals at
`VotingPeriodEnd + MaxExecutionPeriod`.

* ProposalByVotingPeriodEnd:
  `0x13 | sdk.FormatTimeBytes(proposal.VotingPeriodEnd) | BigEndian(ProposalId) -> []byte()`.

## Vote

* Vote: `0x40 | BigEndian(ProposalId) | []byte(voter.Address) -> ProtocolBuffer(Vote)`.

## Grant

Grants are identified by combining grantee address and `Authorization` type
(its target message type URL). Hence we only allow one grant for the (grantee,
Authorization) tuple.

* Grant: `0x20 | len(grant.Grantee) (1 byte) | []byte(grant.Grantee) | []byte(grant.Authorization.MsgTypeURL()) -> ProtocolBuffer(Authorization)`

# Msg Service

## Msg/UpdateParams

The `MsgUpdateParams` can be used to update the parameters of `foundation`.

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/proto/lbm/foundation/v1/tx.proto#L62-L71

It's expected to fail if:

* the authority is not the module's authority.
* the parameters introduces any new foundation-specific features.

## Msg/UpdateDecisionPolicy

The `MsgUpdateDecisionPolicy` can be used to update the decision policy.

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/proto/lbm/foundation/v1/tx.proto#L110-L117

It's expected to fail if:

* the authority is not the module's authority.
* the new decision policy's `Validate()` method doesn't pass.

## Msg/UpdateMembers

Foundation members can be updated with the `MsgUpdateMembers`.

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/proto/lbm/foundation/v1/tx.proto#L97-L105

In the list of `MemberUpdates`, an existing member can be removed by setting
its `remove` flag to true.

It's expected to fail if:

* the authority is not the module's authority.
* if the decision policy's `Validate()` method fails against the updated
  membership.

## Msg/LeaveFoundation

The `MsgLeaveFoundation` allows a foundation member to leave the foundation.

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/proto/lbm/foundation/v1/tx.proto#L205-L209

It's expected to fail if:

* the address is not of a foundation member.
* if the decision policy's `Validate()` method fails against the updated
  membership.

## Msg/SubmitProposal

A new proposal can be created with the `MsgSubmitProposal`, which has a list of
proposers addresses, a list of messages to execute if the proposal is accepted
and some optional metadata.
An optional `Exec` value can be provided to try to execute the proposal
immediately after proposal creation. Proposers signatures are considered as yes
votes in this case.

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/proto/lbm/foundation/v1/tx.proto#L135-L151

It's expected to fail if:

* metadata length is greater than `MaxMetadataLen` config.
* if any of the proposers is not a foundation member.

## Msg/WithdrawProposal

A proposal can be withdrawn using `MsgWithdrawProposal` which has an `address`
(can be either a proposer or the module's authority) and a `proposal_id` (which
has to be withdrawn).

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/proto/lbm/foundation/v1/tx.proto#L159-L166

It's expected to fail if:

* the address is neither the module's authority nor a proposer of the proposal.
* the proposal is already closed or aborted.

## Msg/Vote

A new vote can be created with the `MsgVote`, given a proposal id, a voter
address, a choice (yes, no, veto or abstain) and some optional metadata.
An optional `Exec` value can be provided to try to execute the proposal
immediately after voting.

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/proto/lbm/foundation/v1/tx.proto#L171-L188

It's expected to fail if:

* metadata length is greater than `MaxMetadataLen` config.
* the proposal is not in voting period anymore.

## Msg/Exec

A proposal can be executed with the `MsgExec`.

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/proto/lbm/foundation/v1/tx.proto#L193-L200

The messages that are part of this proposal won't be executed if:

* the proposal has not been accepted by the decision policy.
* the proposal has already been successfully executed.

## Msg/Grant

An authorization grant is created using the `MsgGrant` message.
If there is already a grant for the `(grantee, Authorization)` tuple, then the
new grant overwrites the previous one. To update or extend an existing grant, a
new grant with the same `(grantee, Authorization)` tuple should be created.

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/proto/lbm/foundation/v1/tx.proto#L214-L222

The message handling should fail if:

* the authority is not the module's authority.
* provided `Authorization` is not implemented.
* `Authorization.MsgTypeURL()` is not defined in the router (there is no
  defined handler in the app router to handle that Msg types).

**Note:** Do NOT confuse with that of `x/authz`.

## Msg/Revoke

A grant can be removed with the `MsgRevoke` message.

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/proto/lbm/foundation/v1/tx.proto#L227-L232

The message handling should fail if:

* the authority is not the module's authority.
* provided `MsgTypeUrl` is empty.

## Msg/FundTreasury

Anyone can fund treasury with `MsgFundTreasury`.

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/proto/lbm/foundation/v1/tx.proto#L76-L81

## Msg/WithdrawFromTreasury

The foundation can withdraw coins from the treasury with
`MsgWithdrawFromTreasury`.

+++ https://github.com/line/lbm-sdk/blob/392277a33519d289154e8da27f05f9a6788ab076/proto/lbm/foundation/v1/tx.proto#L86-L92

The message handling should fail if:

* the authority is not the module's authority.
* the address which receives the coins has no authorization of
  `ReceiveFromTreasuryAuthorization`.

## Msg/GovMint

Massive minting is possible through 'MsgGovMint' up to 1 time after the chain is started.

+++ https://github.com/line/lbm-sdk/blob/66988a235a0e01f7a1ee76d719d585ff35f0d176/proto/lbm/foundation/v1/tx.proto#L221-L225

The message handling should fail if:

* the authority is not the module's authority.
* The remaining left count is 0.

# Events

## EventUpdateFoundationParams

`EventUpdateFoundationParams` is an event emitted when the foundation
parameters have been updated.

| Attribute Key | Attribute Value |
|---------------|-----------------|
| params        | {params}        |

## EventUpdateDecisionPolicy

`EventUpdateDecisionPolicy` is an event emitted when the decision policy have
been updated.

| Attribute Key   | Attribute Value  |
|-----------------|------------------|
| decision_policy | {decisionPolicy} |

## EventUpdateMembers

`EventUpdateMembers` is an event emitted when the foundation members have been
updated.

| Attribute Key  | Attribute Value |
|----------------|-----------------|
| member_updates | {members}       |

## EventLeaveFoundation

`EventLeaveFoundation` is an event emitted when a foundation member leaves the
foundation.

| Attribute Key | Attribute Value |
|---------------|-----------------|
| address       | {memberAddress} |

## EventSubmitProposal

`EventSubmitProposal` is an event emitted when a proposal is submitted.

| Attribute Key | Attribute Value |
|---------------|-----------------|
| proposal      | {proposal}      |

## EventWithdrawProposal

`EventWithdrawProposal` is an event emitted when a proposal is withdrawn.

| Attribute Key | Attribute Value |
|---------------|-----------------|
| proposal_id   | {proposalId}    |

## EventVote

`EventVote` is an event emitted when a voter votes on a proposal.

| Attribute Key | Attribute Value |
|---------------|-----------------|
| vote          | {vote}          |

## EventExec

`EventExec` is an event emitted when a proposal is executed.

| Attribute Key | Attribute Value |
|---------------|-----------------|
| proposal_id   | {proposalId}    |
| result        | {result}        |

## EventGrant

`EventGrant` is an event emitted when an authorization is granted to a grantee.

| Attribute Key | Attribute Value  |
|---------------|------------------|
| grantee       | {granteeAddress} |
| authorization | {authorization}  |

## EventRevoke

`EventRevoke` is an event emitted when an authorization is revoked from a
grantee.

| Attribute Key | Attribute Value  |
|---------------|------------------|
| grantee       | {granteeAddress} |
| msg_type_url  | {msgTypeURL}     |

## EventFundTreasury

`EventFundTreasury` is an event emitted when one funds the treasury.

| Attribute Key | Attribute Value |
|---------------|-----------------|
| from          | {fromAddress}   |
| amount        | {amount}        |

## EventWithdrawFromTreasury

`EventWithdrawFromTreasury` is an event emitted when coins are withdrawn from
the treasury.

| Attribute Key | Attribute Value |
|---------------|-----------------|
| to            | {toAddress}     |
| amount        | {amount}        |

## EventGovMint

`EventGovMint` is an event emitted when coins are minted.

| Attribute Key | Attribute Value |
|---------------|-----------------|
| amount        | {amount}        |

# Client

## CLI

A user can query and interact with the `foundation` module using the CLI.

### Query

The `query` commands allow users to query `foundation` state.

```bash
simd query foundation --help
```

#### params

The `params` command allows users to query for the parameters of `foundation`.

```bash
simd query foundation params [flags]
```

Example:

```bash
simd query foundation params
```

Example Output:

```bash
params:
  censored_msg_type_urls:
  - /cosmos.staking.v1beta1.MsgCreateValidator
  - /lbm.foundation.v1.MsgWithdrawFromTreasury
  foundation_tax: "0.200000000000000000"
```

#### foundation-info

The `foundation-info` command allows users to query for the foundation info.

```bash
simd query foundation foundation-info [flags]
```

Example:

```bash
simd query foundation foundation-info
```

Example Output:

```bash
info:
  decision_policy:
    '@type': /lbm.foundation.v1.ThresholdDecisionPolicy
    threshold: "3.000000000000000000"
    windows:
      min_execution_period: 0s
      voting_period: 86400s
  total_weight: "3.000000000000000000"
  version: "1"
```

#### member

The `member` command allows users to query for a foundation member by address.

```bash
simd query foundation member [address] [flags]
```

Example:

```bash
simd query foundation member link1...
```

Example Output:

```bash
member:
  added_at: "0001-01-01T00:00:00Z"
  address: link1...
  metadata: genesis member
```

#### members

The `members` command allows users to query for the foundation members with
pagination flags.

```bash
simd query foundation members [flags]
```

Example:

```bash
simd query foundation members
```

Example Output:

```bash
members:
- added_at: "0001-01-01T00:00:00Z"
  address: link1...
  metadata: genesis member
- added_at: "0001-01-01T00:00:00Z"
  address: link1...
  metadata: genesis member
- added_at: "0001-01-01T00:00:00Z"
  address: link1...
  metadata: genesis member
pagination:
  next_key: null
  total: "3"
```

#### proposal

The `proposal` command allows users to query for proposal by id.

```bash
simd query foundation proposal [id] [flags]
```

Example:

```bash
simd query foundation proposal 1
```

Example Output:

```bash
proposal:
  executor_result: PROPOSAL_EXECUTOR_RESULT_NOT_RUN
  final_tally_result:
    abstain_count: "0.000000000000000000"
    no_count: "0.000000000000000000"
    no_with_veto_count: "0.000000000000000000"
    yes_count: "0.000000000000000000"
  foundation_version: "1"
  id: "1"
  messages:
  - '@type': /lbm.foundation.v1.MsgWithdrawFromTreasury
    authority: link1...
    amount:
    - amount: "1000000000"
      denom: stake
    to: link1...
  metadata: show-me-the-money
  proposers:
  - link1...
  status: PROPOSAL_STATUS_SUBMITTED
  submit_time: "2022-09-19T01:26:38.544943184Z"
  voting_period_end: "2022-09-20T01:26:38.544943184Z"
```

#### proposals

The `proposals` command allows users to query for proposals with pagination
flags.

```bash
simd query foundation proposals [flags]
```

Example:

```bash
simd query foundation proposals
```

Example Output:

```bash
pagination:
  next_key: null
  total: "1"
proposals:
- executor_result: PROPOSAL_EXECUTOR_RESULT_NOT_RUN
  final_tally_result:
    abstain_count: "0.000000000000000000"
    no_count: "0.000000000000000000"
    no_with_veto_count: "0.000000000000000000"
    yes_count: "0.000000000000000000"
  foundation_version: "1"
  id: "1"
  messages:
  - '@type': /lbm.foundation.v1.MsgWithdrawFromTreasury
    authority: link1...
    amount:
    - amount: "1000000000"
      denom: stake
    to: link1...
  metadata: show-me-the-money
  proposers:
  - link1...
  status: PROPOSAL_STATUS_SUBMITTED
  submit_time: "2022-09-19T01:26:38.544943184Z"
  voting_period_end: "2022-09-20T01:26:38.544943184Z"
```

#### vote

The `vote` command allows users to query for vote by proposal id and voter
account address.

```bash
simd query foundation vote [proposal-id] [voter] [flags]
```

Example:

```bash
simd query foundation vote 1 link1...
```

Example Output:

```bash
vote:
  metadata: nope
  option: VOTE_OPTION_NO
  proposal_id: "1"
  submit_time: "2022-09-19T01:35:30.920689570Z"
  voter: link1...
```

#### votes

The `votes` command allows users to query for votes by proposal id with
pagination flags.

```bash
simd query foundation votes [proposal-id] [flags]
```

Example:

```bash
simd query foundation votes 1
```

Example Output:

```bash
pagination:
  next_key: null
  total: "1"
votes:
- metadata: nope
  option: VOTE_OPTION_NO
  proposal_id: "1"
  submit_time: "2022-09-19T01:35:30.920689570Z"
  voter: link1...
```

#### tally

The `tally` command allows users to query for the tally in progress by its
proposal id.

```bash
simd query foundation tally [proposal-id] [flags]
```

Example:

```bash
simd query foundation tally 1
```

Example Output:

```bash
tally:
  abstain_count: "0.000000000000000000"
  no_count: "1.000000000000000000"
  no_with_veto_count: "0.000000000000000000"
  yes_count: "0.000000000000000000"
```

#### grants

The `grants` command allows users to query grants for a grantee. If the message
type URL is set, it selects grants only for that message type.

```bash
simd query foundation grants [grantee] [msg-type-url]? [flags]
```

Example:

```bash
simd query foundation grants link1... /lbm.foundation.v1.MsgWithdrawFromTreasury
```

Example Output:

```bash
authorizations:
- '@type': /lbm.foundation.v1.ReceiveFromTreasuryAuthorization
pagination: null
```

#### treasury

The `treasury` command allows users to query for the foundation treasury.

```bash
simd query foundation treasury [flags]
```

Example:

```bash
simd query foundation treasury
```

Example Output:

```bash
amount:
- amount: "1000000000000.000000000000000000"
  denom: stake
```

### Transactions

The `tx` commands allow users to interact with the `foundation` module.

```bash
simd tx foundation --help
```

**Note:** Some commands must be signed by the module's authority, which means
you cannot broadcast the message directly. The use of those commands is to make
it easier to generate the messages by end users.

#### update-params

The `update-params` command allows users to update the foundation's parameters.

```bash
simd tx foundation update-params [authority] [params-json] [flags]
```

Example:

```bash
simd tx foundation update-params link1... \
    '{
       "foundation_tax": "0.1",
       "censored_msg_type_urls": [
         "/cosmos.staking.v1beta1.MsgCreateValidator",
         "/lbm.foundation.v1.MsgWithdrawFromTreasury"
       ]
     }'
```

**Note:** The signer is the module's authority.

#### update-members

The `update-members` command allows users to update the foundation's members.

```bash
simd tx foundation update-members [authority] [members-json] [flags]
```

Example:

```bash
simd tx foundation update-members link1... \
    '[
       {
         "address": "link1...",
         "metadata": "some new metadata"
       },
       {
         "address": "link1...",
         "remove": true,
       }
     ]'
```

**Note:** The signer is the module's authority.

#### update-decision-policy

The `update-decision-policy` command allows users to update the foundation's
decision policy.

```bash
simd tx foundation update-decision-policy [authority] [decision-policy-json] [flags]
```

Example:

```bash
simd tx foundation update-decision-policy link1... \
    '{
       "@type": "/lbm.foundation.v1.ThresholdDecisionPolicy",
       "threshold": "4",
       "windows": {
         "voting_period": "1h",
         "min_execution_period": "0s"
       }
     }'
```

**Note:** The signer is the module's authority.

#### submit-proposal

The `submit-proposal` command allows users to submit a new proposal.

```bash
simd tx foundation submit-proposal [metadata] [proposers-json] [proposers-json] [messages-json] [flags]
```

Example:

```bash
simd tx foundation submit-proposal show-me-the-money \
    '[
       "link1...",
       "link1..."
     ]' \
    '[
       {
         "@type": "/lbm.foundation.v1.MsgWithdrawFromTreasury",
         "authority": "link1...",
         "to": "link1...",
         "amount": [
           {
             "denom": "stake",
             "amount": "10000000000"
           }
         ]
       }
     ]'
```

#### withdraw-proposal

The `withdraw-proposal` command allows users to withdraw a proposal.

```bash
simd tx foundation withdraw-proposal [proposal-id] [authority-or-proposer] [flags]
```

Example:

```bash
simd tx foundation withdraw-proposal 1 link1...
```

#### vote

The `vote` command allows users to vote on a proposal.

```bash
simd tx foundation vote [proposal-id] [voter] [option] [metadata] [flags]
```

Example:

```bash
simd tx foundation vote 1 link1... VOTE_OPTION_NO nope
```

#### exec

The `exec` command allows users to execute a proposal.

```bash
simd tx foundation exec [proposal-id] [flags]
```

Example:

```bash
simd tx foundation exec 1
```

#### leave-foundation

The `leave-foundation` command allows foundation member to leave the
foundation.

```bash
simd tx foundation leave-foundation [address] [flags]
```

Example:

```bash
simd tx foundation leave-foundation link1...
```

#### grant

The `grant` command allows users to grant an authorization to a grantee.

```bash
simd tx foundation grant [authority] [grantee] [authorization-json] [flags]
```

Example:

```bash
simd tx foundation grant link1.. link1... \
    '{
       "@type": "/lbm.foundation.v1.ReceiveFromTreasuryAuthorization",
     }'
```

**Note:** The signer is the module's authority.

#### revoke

The `revoke` command allows users to revoke an authorization from a grantee.

```bash
simd tx foundation revoke [authority] [grantee] [msg-type-url] [flags]
```

Example:

```bash
simd tx foundation revoke link1.. link1... /lbm.foundation.v1.MsgWithdrawFromTreasury
```

**Note:** The signer is the module's authority.

#### fund-treasury

The `fund-treasury` command allows users to fund the foundation treasury.

```bash
simd tx foundation fund-treasury [from] [amount] [flags]
```

Example:

```bash
simd tx foundation fund-treasury link1.. 1000stake
```

#### withdraw-from-treasury

The `withdraw-from-treasury` command allows users to withdraw coins from the
foundation treasury.

```bash
simd tx foundation withdraw-from-treasury [authority] [to] [amount] [flags]
```

Example:

```bash
simd tx foundation withdraw-from-treasury link1.. link1... 1000stake
```

**Note:** The signer is the module's authority.

## gRPC

A user can query the `foundation` module using gRPC endpoints.

```bash
grpcurl -plaintext \
    localhost:9090 list lbm.foundation.v1.Query
```

### Params

The `Params` endpoint allows users to query for the parameters of `foundation`.

```bash
lbm.foundation.v1.Query/Params
```

Example:

```bash
grpcurl -plaintext \
    localhost:9090 lbm.foundation.v1.Query/Params
```

Example Output:

```bash
{
  "params": {
    "foundationTax": "200000000000000000"
    "censoredMsgTypeUrls": [
      "/cosmos.staking.v1beta1.MsgCreateValidator",
      "/lbm.foundation.v1.MsgWithdrawFromTreasury"
    ]
  }
}
```

### FoundationInfo

The `FoundationInfo` endpoint allows users to query for the foundation info.

```bash
lbm.foundation.v1.Query/FoundationInfo
```

Example:

```bash
grpcurl -plaintext \
    localhost:9090 lbm.foundation.v1.Query/FoundationInfo
```

Example Output:

```bash
{
  "info": {
    "version": "1",
    "totalWeight": "3000000000000000000",
    "decisionPolicy": {"@type":"/lbm.foundation.v1.ThresholdDecisionPolicy","threshold":"3000000000000000000","windows":{"votingPeriod":"86400s","minExecutionPeriod":"0s"}}
  }
}
```

### Member

The `Member` endpoint allows users to query for a foundation member by address.

```bash
lbm.foundation.v1.Query/Member
```

Example:

```bash
grpcurl -plaintext \
    -d '{"address": "link1..."}'
    localhost:9090 lbm.foundation.v1.Query/Member
```

Example Output:

```bash
{
  "member": {
    "address": "link1...",
    "metadata": "genesis member",
    "addedAt": "0001-01-01T00:00:00Z"
  }
}
```

### Members

The `Members` endpoint allows users to query for the foundation members with
pagination flags.

```bash
lbm.foundation.v1.Query/Members
```

Example:

```bash
grpcurl -plaintext \
    localhost:9090 lbm.foundation.v1.Query/Members
```

Example Output:

```bash
{
  "members": [
    {
      "address": "link1...",
      "metadata": "genesis member",
      "addedAt": "0001-01-01T00:00:00Z"
    },
    {
      "address": "link1...",
      "metadata": "genesis member",
      "addedAt": "0001-01-01T00:00:00Z"
    },
    {
      "address": "link1...",
      "metadata": "genesis member",
      "addedAt": "0001-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "total": "3"
  }
}
```

### Proposal

The `Proposal` endpoint allows users to query for proposal by id.

```bash
lbm.foundation.v1.Query/Proposal
```

Example:

```bash
grpcurl -plaintext \
    -d '{"proposal_id": "1"}' \
    localhost:9090 lbm.foundation.v1.Query/Proposal
```

Example Output:

```bash
{
  "proposal": {
    "id": "1",
    "metadata": "show-me-the-money",
    "proposers": [
      "link1..."
    ],
    "submitTime": "2022-09-19T01:26:38.544943184Z",
    "foundationVersion": "1",
    "status": "PROPOSAL_STATUS_SUBMITTED",
    "finalTallyResult": {
      "yesCount": "0",
      "abstainCount": "0",
      "noCount": "0",
      "noWithVetoCount": "0"
    },
    "votingPeriodEnd": "2022-09-20T01:26:38.544943184Z",
    "executorResult": "PROPOSAL_EXECUTOR_RESULT_NOT_RUN",
    "messages": [
      {"@type":"/lbm.foundation.v1.MsgWithdrawFromTreasury","authority":"link1...","amount":[{"denom":"stake","amount":"1000000000"}],"to":"link1..."}
    ]
  }
}
```

### Proposals

The `Proposals` endpoint allows users to query for proposals with pagination
flags.

```bash
lbm.foundation.v1.Query/Proposals
```

Example:

```bash
grpcurl -plaintext \
    localhost:9090 lbm.foundation.v1.Query/Proposals
```

Example Output:

```bash
{
  "proposals": [
    {
      "id": "1",
      "metadata": "show-me-the-money",
      "proposers": [
        "link1..."
      ],
      "submitTime": "2022-09-19T01:26:38.544943184Z",
      "foundationVersion": "1",
      "status": "PROPOSAL_STATUS_SUBMITTED",
      "finalTallyResult": {
        "yesCount": "0",
        "abstainCount": "0",
        "noCount": "0",
        "noWithVetoCount": "0"
      },
      "votingPeriodEnd": "2022-09-20T01:26:38.544943184Z",
      "executorResult": "PROPOSAL_EXECUTOR_RESULT_NOT_RUN",
      "messages": [
        {"@type":"/lbm.foundation.v1.MsgWithdrawFromTreasury","authority":"link1...","amount":[{"denom":"stake","amount":"1000000000"}],"to":"link1..."}
      ]
    }
  ],
  "pagination": {
    "total": "1"
  }
}
```

### Vote

The `Vote` endpoint allows users to query for vote by proposal id and voter account address.

```bash
lbm.foundation.v1.Query/Vote
```

Example:

```bash
grpcurl -plaintext \
    -d '{"proposal_id": "1", "voter": "link1..."}' \
    localhost:9090 lbm.foundation.v1.Query/Vote
```

Example Output:

```bash
{
  "vote": {
    "proposalId": "1",
    "voter": "link1...",
    "option": "VOTE_OPTION_NO",
    "metadata": "nope",
    "submitTime": "2022-09-19T01:35:30.920689570Z"
  }
}
```

### Votes

The `Votes` endpoint allows users to query for votes by proposal id with
pagination flags.

```bash
lbm.foundation.v1.Query/Votes
```

Example:

```bash
grpcurl -plaintext \
    -d '{"proposal_id": "1"}' \
    localhost:9090 lbm.foundation.v1.Query/Votes
```

Example Output:

```bash
{
  "votes": [
    {
      "proposalId": "1",
      "voter": "link1...",
      "option": "VOTE_OPTION_NO",
      "metadata": "nope",
      "submitTime": "2022-09-19T01:35:30.920689570Z"
    }
  ],
  "pagination": {
    "total": "1"
  }
}
```

### TallyResult

The `TallyResult` endpoint allows users to query for the tally in progress by
its proposal id.

```bash
lbm.foundation.v1.Query/Vote
```

Example:

```bash
grpcurl -plaintext \
    -d '{"proposal_id": "1"}' \
    localhost:9090 lbm.foundation.v1.Query/TallyResult
```

Example Output:

```bash
{
  "tally": {
    "yesCount": "0",
    "abstainCount": "0",
    "noCount": "1000000000000000000",
    "noWithVetoCount": "0"
  }
}
```

### Grants

The `Grants` endpoint allows users to query grants for a grantee. If the
message type URL is set, it selects grants only for that message type.

```bash
lbm.foundation.v1.Query/Grants
```

Example:

```bash
grpcurl -plaintext \
    -d '{"grantee": "link1...", "msg_type_url": "/lbm.foundation.v1.MsgWithdrawFromTreasury"}' \
    localhost:9090 lbm.foundation.v1.Query/Grants
```

Example Output:

```bash
{
  "authorizations": [
    {"@type":"/lbm.foundation.v1.ReceiveFromTreasuryAuthorization"}
  ]
}
```

### Treasury

The `Treasury` endpoint allows users to query for the foundation treasury.

```bash
lbm.foundation.v1.Query/Treasury
```

Example:

```bash
grpcurl -plaintext \
    localhost:9090 lbm.foundation.v1.Query/Treasury
```

Example Output:

```bash
{
  "amount": [
    {
      "denom": "stake",
      "amount": "1000000000000000000000000000000"
    }
  ]
}
```
