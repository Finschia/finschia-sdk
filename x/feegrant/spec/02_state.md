<!--
order: 2
-->

# State

## FeeAllowance

Fee Allowances are identified by combining `Grantee` (the account address of fee allowance grantee) with the `Granter` (the account address of fee allowance granter).

<<<<<<< HEAD
Fee allowance grants are stored in the state as follows:

- Grant: `0x00 | grantee_addr_len (1 byte) | grantee_addr_bytes |  granter_addr_len (1 byte) | granter_addr_bytes -> ProtocolBuffer(Grant)`

+++ https://github.com/line/lbm-sdk/blob/691032b8be0f7539ec99f8882caecefc51f33d1f/x/feegrant/feegrant.pb.go#L221-L229
=======
Fee allowances are stored in the state as follows:

- FeeAllowance: `0x00 | grantee_addr_len (1 byte) | grantee_addr_bytes |  granter_addr_len (1 byte) | granter_addr_bytes -> ProtocolBuffer(FeeAllowance)`

+++ https://github.com/cosmos/cosmos-sdk/blob/d97e7907f176777ed8a464006d360bb3e1a223e4/x/feegrant/types/feegrant.pb.go#L358-L363
>>>>>>> 0af248b95 (Add specs for feegrant (#8496))
