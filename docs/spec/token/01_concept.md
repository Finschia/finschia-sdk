# Concept

## Rule for defining symbols

a symbol is composed of alphanumeric.
```regexp
^[a-z][a-z0-9]{2,15}$
```

### Elements of symbols

#### TCK: Ticker
```regexp
[a-z][a-z0-9]{2,4}
```
* Ex) 
  * link, btc, eth, xrp

#### TID: TokenID
```regexp
[a-z0-9]{8}
```
* Ex)
  * 00000000, 00000001
  * sword001, sword002
  * contract0, contract1

#### AAS: Account Address Suffix
```regexp
[a-z0-9]{3}
```
* Ex)
  * jcq (link1qfm2yuj0agtaz76l4455s7qd84nuva0afv2jcq)
  * a22 (link1usmf7gja2qqrs8pnh30jwte0u65k9e6mpjea22)
  * xwr (link1prdg0gjym2q8yzwpuywwn6u9q4lf94sj75axwr)

### Define Rules according the length of the symbols
```
                  user                                                           
       reserved   defined           user defined                                        
       token      token             collective token                                    
       <------> <------>                <------>                                  
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+                                 
|01|02|03|04|05|06|07|08|09|10|11|12|13|14|15|16|                                 
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+                                 
+--------------+                                                               
|       TCK    |                                                               
+--------------+                                                               
|Reserved Token|                                                               
+--------------+                                                               
|link, btc, eth|                                                               
+--------------+                                                               
+--------------+--------+                                                      
|       TCK    |  AAS   |                                                      
+--------------+--------+                                                      
|  User-Defined Token   |                                                      
+-----------------------+                                                      
|   conyjcq, browna22   |                                                      
+-----------------------+                                                      
+--------------+--------+-----------------------+                                 
|       TCK    |  AAS   |        TID            |                                 
+--------------+--------+-----------------------+                                 
|User-Defined Collection|        TID            |                                 
+-----------------------+-----------------------+                                 
|            Collective Token                   |                                 
+-----------------------------------------------+                                 
|     conyjcq00000001, browna22tid00003         |                                 
+-----------------------------------------------+                                 
```
#### Reserved Token
```regexp
^{TCK}$
^[a-z][a-z0-9]{2,4}$
```
A symbol with 3 ~ 5 length is reserved.
An account user cannot occupy any permissions to issue or mint Token with this symbol.

#### User-Defined Token
```regexp
^{TCK}{AAS}$
^[a-z][a-z0-9]{2,4}[a-z0-9]{3}$
```
A Token with 6 ~ 8 length symbols can be issued by an account.
An account can occupy permissions to issue of this token if there is no one occupied the token yet.
The suffix(AAS) should be the same as the last 3 bytes of the address of the account which sends the transaction.

#### User-Defined Collection 
```regexp
^{TCK}{AAS}$
^[a-z][a-z0-9]{2,4}[a-z0-9]{3}$
```
When an account occupies a user-defined token for the symbol, **the same symbol for the user-defined collection is also occupied to the account.**
And the account can issue a user-defined collective token for the collection.

#### User-Defined Collective Token
```regexp
^{TCK}{AAS}{TID}$
^[a-z][a-z0-9]{2,4}[a-z0-9]{3}[a-z0-9]{8}$
```
An account can issue collective token for the collection name.

## Fungible vs. Non Fungible Token
Conceptually, the non-fungible token is a token that exists only one in the world and will not be minted more.
```go
type Token struct {
	Name     string  
	Symbol   string  
	Mintable bool    
	Decimals sdk.Int 
	TokenURI string  
}
```
Here is internal structure for the Token.
The non-fungible token expressed with the Token, where Mintable is false and Decimals is 0.
It is guaranteed that when the token is issued with `total amount=1`, `mintable=false`, `decimals=0`, the token is unique in the world and undividable

## Permissions
### Issue permission
* **reserved token** Initially, nobody has this permission. But may be granted by the governance or huge fee. **Not yet decided**
* **user-defined token** An account can issue relevant(AAS) token. And can grant the permission to others or revoke itself
* **collective-token** An account can issue relevant(AAS) collective token. And can grant the permission to others or revoke itself
### Mint Permission
This permission is granted to the account which issued the token.
* **reserved token**
* **user-defined token**
* **collective-token**
