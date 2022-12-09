<!--
order: 4
-->

# BeginBlock

## Liveness Tracking

At the beginning of each block, we update the `ValidatorSigningInfo` for each
voter and check if they've crossed below the liveness threshold over a
sliding window. This sliding window is defined by `SignedBlocksWindow` and the
index in this window is determined by `VoterSetCounter` found in the voter's
`ValidatorSigningInfo`. For each vote processed, the `VoterSetCounter` is incremented
regardless if the voter signed or not. Once the index is determined, the
`MissedBlocksBitArray` and `MissedBlocksCounter` are updated accordingly.

Finally, in order to determine if a voter crosses below the liveness threshold,
we fetch the maximum number of blocks missed, `maxMissed`, which is
`SignedBlocksWindow - (MinSignedPerWindow * SignedBlocksWindow)` and the minimum
height at which we can determine liveness, `minVoterSetCount`. If the voter set counter is
greater than `minVoterSetCount` and the voter's `MissedBlocksCounter` is greater than
`maxMissed`, they will be slashed by `SlashFractionDowntime`, will be jailed
for `DowntimeJailDuration`, and have the following values reset:
`MissedBlocksBitArray` and `MissedBlocksCounter`.

**Note**: Liveness slashes do **NOT** lead to a tombstombing.

```go
height := ctx.BlockHeight()

for _, voteInfo := range req.LastCommitInfo.getVotes() {
  // fetch the validator public key
  consAddr := sdk.BytesToConsAddress(voteInfo.Validator.Address)
  if _, err := k.GetPubkey(ctx, addr); err != nil {
    panic(fmt.Sprintf("Validator consensus-address %s not found", consAddr))
  }

  // fetch signing info
  signInfo, found := k.GetValidatorSigningInfo(ctx, consAddr)
  if !found {
    panic(fmt.Sprintf("Expected signing info for validator %s but not found", consAddr))
  }

  // this is a relative index, so it counts blocks the validator *should* have signed
  // will use the 0-value default signing info if not present, except for the beginning
  voterSetCounter := signInfo.VoterSetCounter
  signInfo.VoterSetCounter++
  index := voterSetCounter % k.SignedBlocksWindow(ctx)

  // Update signed block bit array & counter
  // This counter just tracks the sum of the bit array
  // That way we avoid needing to read/write the whole array each time
  previous := k.GetValidatorMissedBlockBitArray(ctx, consAddr, index)
  missed := !voteInfo.SignedLastBlock
  switch {
  case !previous && missed:
    // Array value has changed from not missed to missed, increment counter
    k.SetValidatorMissedBlockBitArray(ctx, consAddr, index, true)
    signInfo.MissedBlocksCounter++
  case previous && !missed:
    // Array value has changed from missed to not missed, decrement counter
    k.SetValidatorMissedBlockBitArray(ctx, consAddr, index, false)
    signInfo.MissedBlocksCounter--
  default:
    // Array value at this index has not changed, no need to update counter
  }

  minSignedPerWindow := k.MinSignedPerWindow(ctx)

  if missed {
    // emit events...
  }

  minVoterSetCount := k.SignedBlocksWindow(ctx)
  maxMissed := k.SignedBlocksWindow(ctx) - minSignedPerWindow

  // if we have joined enough times to voter set and the validator has missed too many blocks, punish them
  if voterSetCounter >= minVoterSetCount && signInfo.MissedBlocksCounter > maxMissed {
    validator := k.sk.ValidatorByConsAddr(ctx, consAddr)
    if validator != nil && !validator.IsJailed() {
      // Downtime confirmed: slash and jail the validator
      // We need to retrieve the stake distribution which signed the block, so we subtract ValidatorUpdateDelay from the evidence height,
      // and subtract an additional 1 since this is the LastCommit.
      // Note that this *can* result in a negative "distributionHeight" up to -ValidatorUpdateDelay-1,
      // i.e. at the end of the pre-genesis block (none) = at the beginning of the genesis block.
      // That's fine since this is just used to filter unbonding delegations & redelegations.
      distributionHeight := height - sdk.ValidatorUpdateDelay - 1

      // emit events...

      k.sk.Slash(ctx, consAddr, distributionHeight, voteInfo.Validator.Power, k.SlashFractionDowntime(ctx))
      k.sk.Jail(ctx, consAddr)

      signInfo.JailedUntil = ctx.BlockHeader().Time.Add(k.DowntimeJailDuration(ctx))

      // We need to reset the counter & array so that the validator won't be immediately slashed for downtime upon rebonding.
      signInfo.MissedBlocksCounter = 0
      signInfo.VoterSetCounter = 0
      k.clearValidatorMissedBlockBitArray(ctx, consAddr)

      // log events...
    } else {
      // validator was (a) not found or (b) already jailed so we do not slash

      // log events...
    }
  }

  // Set the updated signing info
  k.SetValidatorSigningInfo(ctx, consAddr, signInfo)
}
```
