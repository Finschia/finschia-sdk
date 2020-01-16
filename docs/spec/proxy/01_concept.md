# Concept

* The Proxy module allows transferring coins on behalf of the coin owner.
* The coin owner may approve/disapprove coins for a proxy
* The proxy may transfer coins on behalf of the owner, within the approved coin amounts.

## Example

### Initial Balances
* `rinahTheOnBehalfOf` : 10 link
* `tinaTheProxy` : 5 link
* `evelynTheReceiver` : 1 link

### Approve Coins 
`rinahTheOnBehalfOf` approves 5 link for `tinaTheProxy`
```
linkcli tx proxy approve $(linkcli keys show tina -a) $(linkcli keys show rinah -a) link 5 
```
* `rinahTheOnBehalfOf`'s sign is required
* No changes on balance

### Send Coins From I
`tinaTheProxy` sends 2 link to `evelynTheReceiver` on behalf of `rinahTheOnBehalfOf`
```
linkcli tx proxy approve $(linkcli keys show tina -a) $(linkcli keys show rinah -a) $(linkcli keys show evelyn -a) link 2 
```
* `tinaTheProxy`'s sign is required
* `rinahTheOnBehalfOf` : 8 link
* `tinaTheProxy` : 5 link
* `evelynTheReceiver` : 3 link

### Disapprove Coins 
`rinahTheOnBehalfOf` disapprove 1 link from `tinaTheProxy`
```
linkcli tx proxy disapprove $(linkcli keys show tina -a) $(linkcli keys show rinah -a) link 1 
```
* `rinahTheOnBehalfOf`'s sign is required
* No changes on balance

### Send Coins From II 
`tinaTheProxy` sends 2 link to `evelynTheReceiver` on behalf of `rinahTheOnBehalfOf`
```
linkcli tx proxy approve $(linkcli keys show tina -a) $(linkcli keys show rinah -a) $(linkcli keys show evelyn -a) link 2 
```
* `tinaTheProxy`'s sign is required
* `rinahTheOnBehalfOf` : 6 link
* `tinaTheProxy` : 5 link
* `evelynTheReceiver` : 5 link