# Events

The proxy module emits the following events:

### MsgProxyApproveCoins

| Type                   | Attribute Key               | Attribute Value        |
|------------------------|-----------------------------|------------------------|
| message                | module                      | proxy                  |
| message                | sender                      | {on_behalf_of_address} |
| message                | action                      | proxy_approve_coins    |
| proxy_approve_coins    | proxy_address               | {proxy_address}        |
| proxy_approve_coins    | proxy_on_behalf_of_address  | {on_behalf_of_address} |
| proxy_approve_coins    | proxy_denom                 | {denom}                |
| proxy_approve_coins    | proxy_amount                | {amount}               |

### MsgProxyDisapproveCoins

| Type                   | Attribute Key               | Attribute Value        |
|------------------------|-----------------------------|------------------------|
| message                | module                      | proxy                  |
| message                | sender                      | {on_behalf_of_address} |
| message                | action                      | proxy_disapprove_coins |
| proxy_disapprove_coins | proxy_address               | {proxy_address}        |
| proxy_disapprove_coins | proxy_on_behalf_of_address  | {on_behalf_of_address} |
| proxy_disapprove_coins | proxy_denom                 | {denom}                |
| proxy_disapprove_coins | proxy_amount                | {amount}               |


### MsgProxySendCoinsFrom

| Type                   | Attribute Key              | Attribute Value        |
|------------------------|----------------------------|------------------------|
| message                | module                     | proxy                  |
| message                | sender                     | {proxy_address}        |
| message                | action                     | proxy_send_coins_from  |
| proxy_send_coins_from  | proxy_address              | {proxy_address}        |
| proxy_send_coins_from  | proxy_on_behalf_of_address | {on_behalf_of_address} |
| proxy_send_coins_from  | proxy_denom                | {denom}                |
| proxy_send_coins_from  | proxy_amount               | {amount}               |
| proxy_send_coins_from  | proxy_to_address           | {to_address}           |