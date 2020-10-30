# Escrow v0.6.0
This is the wasm compiled from https://github.com/CosmWasm/cosmwasm-examples/tree/escrow-0.6.0/escrow with

```
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/code/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/rust-optimizer:0.9.0
```
