<!--
order: 1
-->

# Concepts

## Authorization

The consortium module is designed to contain the authorization information. The other modules may deny its message based on the information of the consortium. As of now, the following modules are using the information:

- **[Staking Plus](../../stakingplus/spec/README.md)**
    - [Msg/CreateValidator](../../stakingplus/spec/03_messages.md#msgcreatevalidator)

One can update the authorization, via proposals:

- `UpdateValidatorAuthsProposal` to authorize `Msg/CreateValidator`

+++ https://github.com/line/lbm-sdk/blob/v0.44.0-rc0/proto/lbm/consortium/v1/consortium.proto#L31-L40

## Shutdown

One can shutdown the consortium module via `UpdateConsortiumParamsProposal`, setting its `params.enabled` to `false`. This process is irreversible, so one cannot re-enable the module.

+++ https://github.com/line/lbm-sdk/blob/v0.44.0-rc0/proto/lbm/consortium/v1/consortium.proto#L20-L29
+++ https://github.com/line/lbm-sdk/blob/v0.44.0-rc0/proto/lbm/consortium/v1/consortium.proto#L9-L12
