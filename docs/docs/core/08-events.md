---
sidebar_position: 1
---
# Events

:::note Synopsis
`Event`s are objects that contain information about the execution of the application. They are mainly used by service providers like block explorers and wallet to track the execution of various messages and index transactions.
:::

## Default Events

There are a few events that are automatically emitted for all messages, directly from `baseapp`.

* `message.action`: The name of the message type. Exactly one exists for each transaction message, in the first ABCI `Event` of each transaction message.
* `message.sender`: The address of the message signer. Exactly one exists for each transaction message, in the first ABCI `Event` of each transaction message.
* `message.module`: The name of the module that emitted the message.

:::tip
The module name is assumed by `baseapp` to be the second element of the message route: `"cosmos.bank.v1beta1.MsgSend" -> "bank"`.
In case a module does not follow the standard message path, (e.g. IBC), it is advised to keep emitting the module name event.
`Baseapp` only emits that event if the module have not already done so.
:::
