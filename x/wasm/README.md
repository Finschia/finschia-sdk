<!-- TOC -->
* [Events](#events)
  * [Message Events](#message-events)
    * [MsgStoreCode](#msgstorecode)
    * [MsgInstantiateContract](#msginstantiatecontract)
    * [MsgStoreCodeAndInstantiateContract](#msgstorecodeandinstantiatecontract)
    * [MsgExecuteContract](#msgexecutecontract)
    * [MsgMigrateContract](#msgmigratecontract)
    * [MsgUpdateAdmin](#msgupdateadmin)
    * [MsgClearAdmin](#msgclearadmin)
  * [Keeper Events](#keeper-events)
    * [Reply](#reply)
    * [Sudo](#sudo)
    * [PinCode](#pincode)
    * [UnpinCode](#unpincode)
  * [Proposal Events](#proposal-events)
<!-- TOC -->

# Events

## Message Events

### MsgStoreCode
| Type       | Attribute Key | Attribute Value         | Note |
|------------|---------------|-------------------------|------|
| message    | module        | wasm                    |      |
| message    | sender        | {senderAddress}         |      |
| store_code | code_id       | {contractCodeID}        |      |
| store_code | feature       | {WasmvmRequiredFeature} |      |

### MsgInstantiateContract
| Type                   | Attribute Key                | Attribute Value                | Note                                          |
|------------------------|------------------------------|--------------------------------|-----------------------------------------------|
| message                | module                       | wasm                           |                                               |
| message                | sender                       | {senderAddress}                |                                               |
| instantiate            | code_id                      | {contractCodeID}               |                                               |
| instantiate            | _contract_address            | {contractAddress}              |                                               |
| transfer               | recipient                    | {recipientAddress}             | Only when the fund exists                     |
| transfer               | sender                       | {senderAddress}                | Only when the fund exists                     |                 
| transfer               | amount                       | {amount}                       | Only when the fund exists                     |
| wasm                   | {customContractAttributeKey} | {customContractAttributeValue} | (optional) Defined by wasm contract developer |
| wasm-{customEventType} | {customContractAttributeKey} | {customContractAttributeKey}   | (optional) Defined by wasm contract developer |

### MsgStoreCodeAndInstantiateContract
| Type                   | Attribute Key                | Attribute Value                 | Note                                          |
|------------------------|------------------------------|---------------------------------|-----------------------------------------------|
| message                | module                       | wasm                            |                                               |
| message                | sender                       | {senderAddress}                 |                                               |
| store_code             | code_id                      | {contractCodeID}                |                                               |
| store_code             | feature                      | {WasmvmRequiredFeature}         |                                               |
| instantiate            | code_id                      | {contractCodeID}                |                                               |
| instantiate            | _contract_address            | {contractAddress}               |                                               |
| transfer               | recipient                    | {recipientAddress}              | Only when the fund exists                     |
| transfer               | sender                       | {senderAddress}                 | Only when the fund exists                     |                 
| transfer               | amount                       | {amount}                        | Only when the fund exists                     |
| wasm                   | {customContractAttributeKey} | {customContractAttributeValue}  | (optional) Defined by wasm contract developer |
| wasm-{customEventType} | {customContractAttributeKey} | {customContractAttributeKey}    | (optional) Defined by wasm contract developer |

### MsgExecuteContract
| Type                   | Attribute Key                | Attribute Value                | Note                                          |
|------------------------|------------------------------|:-------------------------------|-----------------------------------------------|
| message                | module                       | wasm                           |                                               |
| message                | sender                       | {senderAddress}                |                                               |
| execute                | _contract_address            | {contractAddress}              |                                               |
| wasm                   | {customContractAttributeKey} | {customContractAttributeValue} | (optional) Defined by wasm contract developer |
| wasm-{customEventType} | {customContractAttributeKey} | {customContractAttributeKey}   | (optional) Defined by wasm contract developer |

### MsgMigrateContract
| Type                   | Attribute Key                | Attribute Value                | Note                                          |
|------------------------|------------------------------|--------------------------------|-----------------------------------------------|
| message                | module                       | wasm                           |                                               |
| message                | sender                       | {senderAddress}                |                                               |
| migrate                | code_id                      | {newCodeID}                    |                                               |
| migrate                | _contract_address            | {contractAddress}              |                                               |
| wasm                   | {customContractAttributeKey} | {customContractAttributeValue} | (optional) Defined by wasm contract developer |
| wasm-{customEventType} | {customContractAttributeKey} | {customContractAttributeKey}   | (optional) Defined by wasm contract developer |

### MsgUpdateAdmin
| Type    | Attribute Key     | Attribute Value   | Note                      |
|---------|-------------------|-------------------|---------------------------|
| message | module            | wasm              |                           |
| message | sender            | {senderAddress}   |                           |

### MsgClearAdmin
| Type    | Attribute Key     | Attribute Value   | Note                      |
|---------|-------------------|-------------------|---------------------------|
| message | module            | wasm              |                           |
| message | sender            | {senderAddress}   |                           |

## Keeper Events
In addition to message events, the wasm keeper will produce events when the following methods are called (or any method which ends up calling them)

### Reply
`reply` is only called from keeper after processing the submessage

| Type  | Attribute Key     | Attribute Value   | Note |
|-------|-------------------|-------------------|------|
| reply | _contract_address | {contractAddress} |      |

### Sudo
`Sudo` allows priviledged access to a contract. This can never be called by an external tx, but only by another native Go module directly.

| Type | Attribute Key     | Attribute Value   | Note |
|------|-------------------|-------------------|------|
| sudo | _contract_address | {contractAddress} |      |

### PinCode
`PinCode` pins the wasm contract in wasmvm cache.

| Type     | Attribute Key | Attribute Value | Note |
|----------|---------------|-----------------|------|
| pin_code | code_id       | {codeID}        |      |

### UnpinCode

| Type       | Attribute Key | Attribute Value | Note |
|------------|---------------|-----------------|------|
| unpin_code | code_id       | {codeID}        |      |

## Proposal Events
If you use wasm proposal, it makes common event like below.

| Type                | Attribute Key | Attribute Value    | Note |
|---------------------|---------------|--------------------|------|
| gov_contract_result | result        | {resultOfProposal} |      |
