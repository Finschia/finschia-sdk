# Events

The bank module emits the following events:

## Handlers

### MsgSend

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| message  | module        | bank               |
| message  | action        | send               |
| message  | sender        | {senderAddress}    |
| transfer | recipient     | {recipientAddress} |
| transfer | amount        | {amount}           |

### MsgMultiSend

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| message  | module        | bank               |
| message  | action        | multisend          |
| message  | sender        | {senderAddress}    |
| transfer | recipient     | {recipientAddress} |
