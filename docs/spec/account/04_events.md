# Events
**Not fully documented yet** 
The account module emits the following events:


### MsgCreateAccount
| Type            | Attribute Key         | Attribute Value   |
|-----------------|-----------------------|-------------------|
| message         | module                | account           |
| message         | sender                | {fromAddress}     | 
| message         | action                | create_account    |
| create_account  | create_account_from   | {fromAddress}     |
| create_account  | create_account_target | {targetAddress}   |

### MsgEmpty    
| Type            | Attribute Key       | Attribute Value     |
|-----------------|---------------------|---------------------|
| message         | module              | account             |
| message         | sender              | {from}              | 
| message         | action              | empty               |
