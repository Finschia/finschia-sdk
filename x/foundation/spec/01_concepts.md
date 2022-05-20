<!--
order: 1
-->

# Concepts

## Authorization

The foundation module is designed to contain the authorization information. The other modules may deny its message based on the information of the foundation. As of now, the following modules are using the information:

- **[Staking Plus](../../stakingplus/spec/README.md)**
    - [Msg/CreateValidator](../../stakingplus/spec/03_messages.md#msgcreatevalidator)

One can update the authorization, via proposals:

- `UpdateValidatorAuthsProposal` to authorize `Msg/CreateValidator`

    +++ https://github.com/line/lbm-sdk/blob/v0.44.0-rc0/proto/lbm/foundation/v1/foundation.proto#L31-L40
    ```go
    // UpdateValidatorAuthsProposal details a proposal to update validator auths on foundation.
    message UpdateValidatorAuthsProposal {
      option (gogoproto.equal)            = false;
      option (gogoproto.goproto_getters)  = false;
      option (gogoproto.goproto_stringer) = false;

      string                 title       = 1;
      string                 description = 2;
      repeated ValidatorAuth auths       = 3 [(gogoproto.moretags) = "yaml:\"auths\""];
    }
    ```

## Disable the module

One can disable the foundation module via `UpdateFoundationParamsProposal`, setting its `params.enabled` to `false`. This process is irreversible, so one cannot re-enable the module.

+++ https://github.com/line/lbm-sdk/blob/v0.44.0-rc0/proto/lbm/foundation/v1/foundation.proto#L20-L29
```go
// UpdateFoundationParamsProposal details a proposal to update params of foundation module.
message UpdateFoundationParamsProposal {
  option (gogoproto.equal)            = false;
  option (gogoproto.goproto_getters)  = false;
  option (gogoproto.goproto_stringer) = false;

  string title       = 1;
  string description = 2;
  Params params      = 3;
}
```

+++ https://github.com/line/lbm-sdk/blob/v0.44.0-rc0/proto/lbm/foundation/v1/foundation.proto#L9-L12
```go
// Params defines the parameters for the foundation module.
message Params {
  bool enabled = 1 [(gogoproto.moretags) = "yaml:\"enabled\""];
}
```
