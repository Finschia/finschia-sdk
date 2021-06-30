<!--
order: 0
title: IBC Overview
parent:
  title: "ibc"
-->

# `ibc`

## Abstract

This specification defines the implementation of the IBC protocol on the LFB SDK, the
changes made to the specification and where to find each specific ICS spec within
the module.

For the general specification please refer to the [Interchain Standards](https://github.com/cosmos/ics).

## Contents

1. **Applications**

    1.1. [Transfer](./../applications/transfer/spec/README.md)
2. **[Core](./../core/spec/README.md)**
3. **Light Clients**

    3.1 [Solo Machine Client](./../light-clients/06-solomachine/spec/README.md)

    3.2 [Tendermint Client](./../light-clients/07-tendermint/spec/README.md)

    3.3 [Localhost Client](./../light-clients/09-localhost/spec/README.md)

## Implementation Details

As stated above, the IBC implementation on the LFB SDK introduces some changes
to the general specification, in order to avoid code duplication and to take
advantage of the SDK architectural components such as the transaction routing
through `Handlers`.

### Interchain Standards reference

The following list is a mapping from each Interchain Standard to their implementation
in the SDK's `x/ibc` module:

* [ICS 002 - Client Semantics](https://github.com/cosmos/ibc/tree/master/spec/core/ics-002-client-semantics): Implemented in [`x/ibc/core/02-client`](github.com/line/lfb-sdk/tree/main/x/ibc/core/02-client)
* [ICS 003 - Connection Semantics](https://github.com/cosmos/ibc/tree/master/spec/core/ics-003-connection-semantics): Implemented in [`x/ibc/core/03-connection`](github.com/line/lfb-sdk/tree/main/x/ibc/core/03-connection)
* [ICS 004 - Channel and Packet Semantics](https://github.com/cosmos/ibc/tree/master/spec/core/ics-004-channel-and-packet-semantics): Implemented in [`x/ibc/core/04-channel`](github.com/line/lfb-sdk/tree/main/x/ibc/core/04-channel)
* [ICS 005 - Port Allocation](https://github.com/cosmos/ics/blob/master/spec/core/ics-005-port-allocation): Implemented in [`x/ibc/core/05-port`](github.com/line/lfb-sdk/tree/main/x/ibc/core/05-port)
* [ICS 006 - Solo Machine Client](https://github.com/cosmos/ics/blob/master/spec/client/ics-006-solo-machine-client): Implemented in [`x/ibc/light-clients/06-solomachine`](https://github.com/line/lfb-sdk/tree/main/x/ibc/light-clients/06-solomachine)
* [ICS 007 - Tendermint Client](https://github.com/cosmos/ics/blob/master/spec/client/ics-007-tendermint-client): Implemented in [`x/ibc/light-clients/07-tendermint`](https://github.com/line/lfb-sdk/tree/main/x/ibc/light-clients/07-tendermint)
* [ICS 009 - Loopback Client](https://github.com/cosmos/ibc/tree/master/spec/client/ics-009-loopback-client):  Implemented in [`x/ibc/light-clients/09-localhost`](https://github.com/line/lfb-sdk/tree/main/x/ibc/light-clients/09-localhost)
* [ICS 018- Relayer Algorithms](https://github.com/cosmos/ibc/tree/master/spec/relayer/ics-018-relayer-algorithms): Not implemented yet.
* [ICS 020 - Fungible Token Transfer](https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer): Implemented in [`x/ibc/applications/transfer`](https://github.com/line/lfb-sdk/tree/main/x/ibc/applications/transfer)
* [ICS 023 - Vector Commitments](https://github.com/cosmos/ibc/tree/master/spec/core/ics-023-vector-commitments): Implemented in [`x/ibc/core/23-commitment`](github.com/line/lfb-sdk/tree/main/x/ibc/core/23-commitment)
* [ICS 024 - Host Requirements](https://github.com/cosmos/ics/tree/master/spec/core/ics-024-host-requirements): Implemented in [`x/ibc/core/24-host`](github.com/line/lfb-sdk/tree/main/x/ibc/core/24-host)
* [ICS 025 - Handler Interface](https://github.com/cosmos/ics/tree/master/spec/core/ics-025-handler-interface): `Handler` interfaces are implemented at the top level in `x/ibc/handler.go`,
which call each ICS submodule's handlers (i.e `x/ibc/*/{XX-ICS}/handler.go`).

### SDK Modules

* [`x/capability`](https://github.com/line/lfb-sdk/tree/main/x/capability): The capability module provides object-capability keys support through scoped keepers in order to authenticate usage of ports or channels. Check [ADR 3 - Dynamic Capability Store](../../../docs/architecture/adr-003-dynamic-capability-store.md) for more details.

## IBC module architecture

For ease of auditing, every Interchain Standard has been developed in its own
package. The development team separated the IBC TAO (Transport, Authentication, Ordering) ICS specifications from the IBC application level
specification. The following tree describes the architecture of the directories that
the `ibc` (TAO) and `ibc-transfer` ([ICS20](https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer)) modules:

```shell
x/ibc
├── applications/
│   └──transfer/
├── core/
│   ├── 02-client/
│   ├── 03-connection/
│   ├── 04-channel/
│   ├── 05-port/
│   ├── 23-commitment/
│   ├── 24-host/
│   ├── client
│   │   └── cli
│   │       └── cli.go
│   ├── keeper
│   │   ├── keeper.go
│   │   └── querier.go
│   ├── types
│   │   ├── errors.go
│   │   └── keys.go
│   ├── handler.go
│   └── module.go
├── light-clients/
│   ├── 06-solomachine/
│   ├── 07-tendermint/
│   └── 09-localhost/
└── testing/
```
