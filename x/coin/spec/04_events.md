# Events

The coin module emits the following events:

## Handlers

### MsgSend

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| message  | module        | coin               |
| message  | action        | send               |
| message  | sender        | {senderAddress}    |
| transfer | recipient     | {recipientAddress} |
| transfer | amount        | {amount}           |

### MsgMultiSend

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| message  | module        | coin               |
| message  | action        | multisend          |
| message  | sender        | {senderAddress}    |
| transfer | recipient     | {recipientAddress} |
