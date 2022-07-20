
# Events of LBM-SDK

Note: the content of this document is based on the actual implementation in version 0.45.6 of the cosmos-sdk. As we do not intend to introduce a new specification, it's a bug of the documentation if you find any differences between this document and the implementation.

Note: these specifications are subject to change.

## Common events

In its process, every message emits an event of type `message` which has an attribute of a key `action`. This attribute has a value of `TypeURL` of the message. `TypeURL` is just a message name of the message, prefixed by `/`. For example, `TypeURL` of a message `MsgSend` in x/bank would be `/cosmos.bank.v1beta1.MsgSend`.

This event MUST proceed to any other events in its process of the message.

## Ordering

The ordering of the events MUST be preserved in `Events` of `TxResponse`.

Note: the field `Logs` is non-deterministic as cosmos-sdk says.
