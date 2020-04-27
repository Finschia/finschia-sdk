# Concept

## What is the Safety Box?

Safety box allows operator to designate roles in order to manage token allocation and issuance within the box. 
This is only allowed for token, not base coin.

## Roles

There are five roles with different permissions. 

1. Owner
    - The owner must be specified while creating a safety box.
    - Only the owner may register/deregister operators. 
    - There is only one owner account and cannot be changed.
    - Ex. `linkcli tx safetybox create [id] [owner_address] [contract_id]`
2. Operator
    - The operator oversees the safety box. Only the owner may register and deregister operators.
    - Only operator(s) may register allocators, issuers and returners.
    - Multiple operators are allowed.
    - Ex. `linkcli tx safetybox role [safety_box_id] [register|deregister] [operator|allocator|issuer|returner] [from_address] [to_address]`
3. Allocator
    - Only allocators may allocate token to the safety box.
    - Ex. `linkcli tx safetybox sendtoken [safety_box_id] allocate [contract_id] [amount] [address]`
    - Only allocators may recall token from the safety box.
    - Ex. `linkcli tx safetybox sendtoken [safety_box_id] recall [contract_id] [amount] [address]`
    - Multiple allocators are allowed.
4. Issuer
    - Only issuers may request issuance to the safety box. 
    - Only issuers may receive issuance from the safety box.
    - Multiple issuers are allowed.
    - Ex. `linkcli tx safetybox sendtoken [safety_box_id] issue [contract_id] [amount] [address] [issuer_address]`
5. Returner
    - Only returners may return token back to the safety box.
    - Multiple returners are allowed.
    - Ex. `linkcli tx safetybox sendtoken [safety_box_id] return [contract_id] [amount] [address]` 
    
## Safety Box ID Rules

No rules yet but duplicated IDs will be rejected. 
