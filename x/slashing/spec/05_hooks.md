<!--
order: 5
-->

# Hooks

In this section we describe the "hooks" - slashing module code that runs when other events happen.

## Validator Bonded

Upon successful first-time bonding of a new validator, we create a new `ValidatorSigningInfo` structure for the
now-bonded validator.

```
onValidatorBonded(address sdk.ValAddress)

  signingInfo, found = GetValidatorSigningInfo(address)
  if !found {
    signingInfo = ValidatorSigningInfo {
      JailedUntil         : time.Unix(0, 0),
      Tombstone           : false,
      MissedBloskCounter  : 0,
	  VoterSetCounter     : 0,
    }
    setValidatorSigningInfo(signingInfo)
  }

  return
```
