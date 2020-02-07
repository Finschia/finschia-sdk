# Link Contract Tests
The link cli contract test in this folder.

Test at once
```
make contract-test
```

Test sequentially
```
make setup-transactions
make dredd-test
make clean-up-dredd-test
```

### Test Structure
The `contract-test` verifies that the actual rest-api on the LCD matches `swagger.yaml` using dredd.

It run the chain and send a each request in the paths of swagger.yaml to compare the expected and actual values.

The procedure for `contract-test` is as follows:
1. Compile `cmd/contract_test_hook`
2. Start link chain
3. Send a transaction to set the state required for the test in `contract_test/testdata/setup.sh`
4. Run `dredd`

The procedure for `dredd` is as follows:
  1. Run linkcli rest server
  2. Run `contract_test_hook` server
  3. Parse the response and create the `expected json`
  4. Parse the parameters and create the request json
  5. Send request to rest server to get `actual json`
  6. Run before-hooks
  7. Compare `expected json` with `actual json`
  8. Run after-hooks

### Notes when update swagger.yaml
- What to do when adding a new message:
  - Add raw json of the message you added to examples.MsgExamples
  - If there is a new signer in the message that does not exist, you need to add a signature to `contract_test/testdata/setup.sh`
- You need to setup that request to succeed through `contract_test/testdata/setup.sh`
    - Messages in `examples.MsgExamples` are sent to the chain. Add a message to `examples.MsgExamples` to set the status.
    - Address and token_id change dynamically with each run, so modify `contract_test/testdata/replace_symbols.sh` to replace the dummy value in the `swagger.yaml` with the actual value.
    - dummy values are defined in `contract_test/testdata/common.sh`
- `dredd` does not support the oneOf, allOf, anyOf swagger syntax now, so in that case add a hook to verify it.