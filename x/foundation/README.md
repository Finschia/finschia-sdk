<!--
order: 0
title: Foundation Overview
parent:
  title: "foundation"
-->

# `x/foundation`

## Abstract

This module provides the functionalities related to the foundation. The foundation can turn off these functionalities irreversibly, through the corresponding proposal. Therefore, the users can ensure that no one can bring back these foundation-specific functionalities.

## Contents

* [Concepts](#concepts)
* [State](#state)
* [Msg Service](#msg-service)
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
* [Client](#client)
    * [CLI](#cli)
    * [gRPC](#grpc)
* [Parameters](#parameters)
