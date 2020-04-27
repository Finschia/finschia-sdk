# Events

The safety box module emits the following events:


### MsgSafetyBoxCreate

| Type              | Attribute Key       | Attribute Value        |
|-------------------|---------------------|------------------------|
| message           | module              | safety_box             |
| message           | sender              | {safety_box_owner}     |
| message           | action              | safety_box_create      |
| safety_box_create | safety_box_id       | {safety_box_id}        |
| safety_box_create | safety_box_owner    | {safety_box_owner}     |
| safety_box_create | safety_box_address  | {safety_box_address}   |
| safety_box_create | contract_id         | {contract_id}          |


### MsgSafetyBoxAllocateToken

| Type                  | Attribute Key                 | Attribute Value           |
|-----------------------|-------------------------------|---------------------------|
| message               | module                        | safety_box                |
| message               | sender                        | {allocator_address}       |
| message               | action                        | safety_box_allocate_token |
| safety_box_send_token | safety_box_id                 | {safety_box_id}           |
| safety_box_send_token | safety_box_allocator_address  | {allocator_address}       |
| safety_box_send_token | safety_box_action             | allocate                  |
| safety_box_send_token | contract_id                   | {contract_id}             |
| safety_box_send_token | amount                        | {amount}                  |


### MsgSafetyBoxRecallToken

| Type                  | Attribute Key                 | Attribute Value         |
|-----------------------|-------------------------------|-------------------------|
| message               | module                        | safety_box              |
| message               | sender                        | {allocator_address}     |
| message               | action                        | safety_box_recall_token |
| safety_box_send_token | safety_box_id                 | {safety_box_id}         |
| safety_box_send_token | safety_box_allocator_address  | {allocator_address}     |
| safety_box_send_token | safety_box_action             | recall                  |
| safety_box_send_token | contract_id                   | {contract_id}           |
| safety_box_send_token | amount                        | {amount}                |


### MsgSafetyBoxIssueToken

| Type                  | Attribute Key                 | Attribute Value        |
|-----------------------|-------------------------------|------------------------|
| message               | module                        | safety_box             |
| message               | sender                        | {from_address}         |
| message               | action                        | safety_box_issue_token |
| safety_box_send_token | safety_box_id                 | {safety_box_id}        |
| safety_box_send_token | safety_box_issue_from_address | {from_address}         |
| safety_box_send_token | safety_box_issue_to_address   | {to_address}           |
| safety_box_send_token | safety_box_action             | issue                  |
| safety_box_send_token | contract_id                   | {contract_id}          |
| safety_box_send_token | amount                        | {amount}               |


### MsgSafetyBoxReturnToken

| Type                  | Attribute Key                 | Attribute Value         |
|-----------------------|-------------------------------|-------------------------|
| message               | module                        | safety_box              |
| message               | sender                        | {returner_address}      |
| message               | action                        | safety_box_return_token |
| safety_box_send_token | safety_box_id                 | {safety_box_id}         |
| safety_box_send_token | safety_box_returner_address   | {returner_address}      |
| safety_box_send_token | safety_box_action             | return                  |
| safety_box_send_token | contract_id                   | {contract_id}           |
| safety_box_send_token | amount                        | {amount}                |


### MsgSafetyBoxRegisterOperator

| Type                  | Attribute Key                 | Attribute Value                      |
|-----------------------|-------------------------------|--------------------------------------|
| message               | module                        | safety_box                           |
| message               | sender                        | {safety_box_owner}                   |
| message               | action                        | safety_box_grant_operator_permission |
| safety_box_permission | safety_box_id                 | {safety_box_id}                      |
| safety_box_permission | safety_box_owner              | {safety_box_owner}                   |
| safety_box_permission | safety_box_target             | {address}                            |
| safety_box_permission | safety_box_grant_permission   | operator                             |


### MsgSafetyBoxDeregisterOperator

| Type                  | Attribute Key                 | Attribute Value                       |
|-----------------------|-------------------------------|---------------------------------------|
| message               | module                        | safety_box                            |
| message               | sender                        | {safety_box_owner}                    |
| message               | action                        | safety_box_revoke_operator_permission |
| safety_box_permission | safety_box_id                 | {safety_box_id}                       |
| safety_box_permission | safety_box_owner              | {safety_box_owner}                    |
| safety_box_permission | safety_box_target             | {address}                             |
| safety_box_permission | safety_box_revoke_permission  | operator                              |


### MsgSafetyBoxRegisterAllocator

| Type                  | Attribute Key                 | Attribute Value                       |
|-----------------------|-------------------------------|---------------------------------------|
| message               | module                        | safety_box                            |
| message               | sender                        | {operator_address}                    |
| message               | action                        | safety_box_grant_allocator_permission |
| safety_box_permission | safety_box_id                 | {safety_box_id}                       |
| safety_box_permission | safety_box_operator           | {operator_address}                    |
| safety_box_permission | safety_box_target             | {address}                             |
| safety_box_permission | safety_box_grant_permission   | allocator                             |


### MsgSafetyBoxDeregisterAllocator

| Type                  | Attribute Key                 | Attribute Value                        |
|-----------------------|-------------------------------|----------------------------------------|
| message               | module                        | safety_box                             |
| message               | sender                        | {operator_address}                     |
| message               | action                        | safety_box_revoke_allocator_permission |
| safety_box_permission | safety_box_id                 | {safety_box_id}                        |
| safety_box_permission | safety_box_operator           | {operator_address}                     |
| safety_box_permission | safety_box_target             | {address}                              |
| safety_box_permission | safety_box_revoke_permission  | allocator                              |


### MsgSafetyBoxRegisterIssuer

| Type                  | Attribute Key                 | Attribute Value                        |
|-----------------------|-------------------------------|----------------------------------------|
| message               | module                        | safety_box                             |
| message               | sender                        | {operator_address}                     |
| message               | action                        | safety_box_grant_issuer_permission     |
| safety_box_permission | safety_box_id                 | {safety_box_id}                        |
| safety_box_permission | safety_box_operator           | {operator_address}                     |
| safety_box_permission | safety_box_target             | {address}                              |
| safety_box_permission | safety_box_grant_permission   | issuer                                 |


### MsgSafetyBoxDeregisterIssuer

| Type                  | Attribute Key                 | Attribute Value                     |
|-----------------------|-------------------------------|-------------------------------------|
| message               | module                        | safety_box                          |
| message               | sender                        | {operator_address}                  |
| message               | action                        | safety_box_revoke_issuer_permission |
| safety_box_permission | safety_box_id                 | {safety_box_id}                     |
| safety_box_permission | safety_box_operator           | {operator_address}                  |
| safety_box_permission | safety_box_target             | {address}                           |
| safety_box_permission | safety_box_revoke_permission  | issuer                              |


### MsgSafetyBoxRegisterReturner

| Type                  | Attribute Key                 | Attribute Value                      |
|-----------------------|-------------------------------|--------------------------------------|
| message               | module                        | safety_box                           |
| message               | sender                        | {operator_address}                   |
| message               | action                        | safety_box_grant_returner_permission |
| safety_box_permission | safety_box_id                 | {safety_box_id}                      |
| safety_box_permission | safety_box_operator           | {operator_address}                   |
| safety_box_permission | safety_box_target             | {address}                            |
| safety_box_permission | safety_box_grant_permission   | returner                             |


### MsgSafetyBoxDeregisterReturner

| Type                  | Attribute Key                 | Attribute Value                       |
|-----------------------|-------------------------------|---------------------------------------|
| message               | module                        | safety_box                            |
| message               | sender                        | {operator_address}                    |
| message               | action                        | safety_box_revoke_returner_permission |
| safety_box_permission | safety_box_id                 | {safety_box_id}                       |
| safety_box_permission | safety_box_operator           | {operator_address}                    |
| safety_box_permission | safety_box_target             | {address}                             |
| safety_box_permission | safety_box_revoke_permission  | returner                              |
