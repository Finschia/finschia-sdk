<!-- TOC -->
# Category
  * [Authz](#authz)
  * [Bank](#bank)
  * [Capability](#capability)
  * [Collection](#collection)
  * [Crisis](#crisis)
  * [Distribution](#distribution)
  * [Evidence](#evidence)
  * [Feegrant](#feegrant)
  * [Foundation](#foundation)
  * [Gov](#gov)
  * [Params](#params)
  * [Slashing](#slashing)
  * [Staking](#staking)
  * [Token](#token)
<!-- TOC -->

## Authz

|Error Name|Codespace|Code|Description|
|:-|:-|:-|:-|
|ErrInvalidExpirationTime|authz|3|expiration time of authorization should be more than current time|

>You can also find detailed information in the following Errors.go files:
  * [authz/errors.go](authz/errors.go)

## Bank

|Error Name|Codespace|Code|Description|
|:-|:-|:-|:-|
|ErrNoInputs|bank|2|no inputs to send transaction|
|ErrNoOutputs|bank|3|no outputs to send transaction|
|ErrInputOutputMismatch|bank|4|sum inputs != sum outputs|
|ErrSendDisabled|bank|5|send transactions are disabled|
|ErrDenomMetadataNotFound|bank|6|client denom metadata not found|
|ErrInvalidKey|bank|7|invalid key|

>You can also find detailed information in the following Errors.go files:
  * [bank/types/errors.go](bank/types/errors.go)

## Capability

|Error Name|Codespace|Code|Description|
|:-|:-|:-|:-|
|ErrInvalidCapabilityName|capability|2|capability name not valid|
|ErrNilCapability|capability|3|provided capability is nil|
|ErrCapabilityTaken|capability|4|capability name already taken|
|ErrOwnerClaimed|capability|5|given owner already claimed capability|
|ErrCapabilityNotOwned|capability|6|capability not owned by module|
|ErrCapabilityNotFound|capability|7|capability not found|
|ErrCapabilityOwnersNotFound|capability|8|owners not found for capability|

>You can also find detailed information in the following Errors.go files:
  * [capability/types/errors.go](capability/types/errors.go)

## Collection

|Error Name|Codespace|Code|Description|
|:-|:-|:-|:-|
|ErrTokenNotExist|collection|2|token symbol, token-id does not exist|
|ErrTokenNotMintable|collection|3|token symbol, token-id is not mintable|
|ErrInvalidTokenName|collection|4|token name should not be empty|
|ErrInvalidTokenID|collection|5|invalid token id|
|ErrInvalidTokenDecimals|collection|6|token decimals should be within the range in 0 ~ 18|
|ErrInvalidIssueFT|collection|7|Issuing token with amount[1], decimals[0], mintable[false] prohibited. Issue nft token instead.|
|ErrInvalidAmount|collection|8|invalid token amount|
|ErrInvalidBaseImgURILength|collection|9|invalid base_img_uri length|
|ErrInvalidNameLength|collection|10|invalid name length|
|ErrInvalidTokenType|collection|11|invalid token type pattern found|
|ErrInvalidTokenIndex|collection|12|invalid token index pattern found|
|ErrCollectionExist|collection|13|collection already exists|
|ErrCollectionNotExist|collection|14|collection does not exists|
|ErrTokenTypeExist|collection|15|token type for contract_id, token-type already exists|
|ErrTokenTypeNotExist|collection|16|token type for contract_id, token-type does not exist|
|ErrTokenTypeFull|collection|17|all token type for contract_id are occupied|
|ErrTokenIndexFull|collection|18|all non-fungible token index for contract_id, token-type are occupied|
|ErrTokenIDFull|collection|19|all fungible token-id for contract_id are occupied|
|ErrTokenNoPermission|collection|20|account does not have the permission|
|ErrTokenAlreadyAChild|collection|21|token is already a child of some other|
|ErrTokenNotAChild|collection|22|token is not a child of some other|
|ErrTokenNotOwnedBy|collection|23|token is being not owned by|
|ErrTokenCannotTransferChildToken|collection|24|cannot transfer a child token|
|ErrTokenNotNFT|collection|25|token is not a NFT|
|ErrCannotAttachToItself|collection|26|cannot attach token to itself|
|ErrCannotAttachToADescendant|collection|27|cannot attach token to a descendant|
|ErrApproverProxySame|collection|28|approver is same with proxy|
|ErrCollectionNotApproved|collection|29|proxy is not approved on the collection|
|ErrCollectionAlreadyApproved|collection|30|proxy is already approved on the collection|
|ErrAccountExist|collection|31|account already exists|
|ErrAccountNotExist|collection|32|account does not exists|
|ErrInsufficientSupply|collection|33|insufficient supply|
|ErrInvalidCoin|collection|34|invalid coin|
|ErrInvalidChangesFieldCount|collection|35|invalid count of field changes|
|ErrEmptyChanges|collection|36|changes is empty|
|ErrInvalidChangesField|collection|37|invalid field of changes|
|ErrTokenIndexWithoutType|collection|38|There is a token index but no token type|
|ErrTokenTypeFTWithoutIndex|collection|39|There is a token type of ft but no token index|
|ErrInsufficientToken|collection|40|insufficient token|
|ErrDuplicateChangesField|collection|41|duplicate field of changes|
|ErrInvalidMetaLength|collection|42|invalid meta length|
|ErrSupplyOverflow|collection|43|supply for collection reached maximum|
|ErrEmptyField|collection|44|required field cannot be empty|
|ErrCompositionTooDeep|collection|45|cannot attach token (composition too deep)|
|ErrCompositionTooWide|collection|46|cannot attach token (composition too wide)|
|ErrBurnNonRootNFT|collection|47|cannot burn non-root NFTs|

>You can also find detailed information in the following Errors.go files:
  * [collection/errors.go](collection/errors.go)

## Crisis

|Error Name|Codespace|Code|Description|
|:-|:-|:-|:-|
|ErrNoSender|crisis|2|sender address is empty|
|ErrUnknownInvariant|crisis|3|unknown invariant|

>You can also find detailed information in the following Errors.go files:
  * [crisis/types/errors.go](crisis/types/errors.go)

## Distribution

|Error Name|Codespace|Code|Description|
|:-|:-|:-|:-|
|ErrEmptyDelegatorAddr|distribution|2|delegator address is empty|
|ErrEmptyWithdrawAddr|distribution|3|withdraw address is empty|
|ErrEmptyValidatorAddr|distribution|4|validator address is empty|
|ErrEmptyDelegationDistInfo|distribution|5|no delegation distribution info|
|ErrNoValidatorDistInfo|distribution|6|no validator distribution info|
|ErrNoValidatorCommission|distribution|7|no validator commission to withdraw|
|ErrSetWithdrawAddrDisabled|distribution|8|set withdraw address disabled|
|ErrBadDistribution|distribution|9|community pool does not have sufficient coins to distribute|
|ErrInvalidProposalAmount|distribution|10|invalid community pool spend proposal amount|
|ErrEmptyProposalRecipient|distribution|11|invalid community pool spend proposal recipient|
|ErrNoValidatorExists|distribution|12|validator does not exist|
|ErrNoDelegationExists|distribution|13|delegation does not exist|

>You can also find detailed information in the following Errors.go files:
  * [distribution/types/errors.go](distribution/types/errors.go)

## Evidence

|Error Name|Codespace|Code|Description|
|:-|:-|:-|:-|
|ErrNoEvidenceHandlerExists|evidence|2|unregistered handler for evidence type|
|ErrInvalidEvidence|evidence|3|invalid evidence|
|ErrNoEvidenceExists|evidence|4|evidence does not exist|
|ErrEvidenceExists|evidence|5|evidence already exists|

>You can also find detailed information in the following Errors.go files:
  * [evidence/types/errors.go](evidence/types/errors.go)

## Feegrant

|Error Name|Codespace|Code|Description|
|:-|:-|:-|:-|
|ErrFeeLimitExceeded|feegrant|2|fee limit exceeded|
|ErrFeeLimitExpired|feegrant|3|fee allowance expired|
|ErrInvalidDuration|feegrant|4|invalid duration|
|ErrNoAllowance|feegrant|5|no allowance|
|ErrNoMessages|feegrant|6|allowed messages are empty|
|ErrMessageNotAllowed|feegrant|7|message not allowed|

>You can also find detailed information in the following Errors.go files:
  * [feegrant/errors.go](feegrant/errors.go)

## Foundation

|Error Name|Codespace|Code|Description|
|:-|:-|:-|:-|

>You can also find detailed information in the following Errors.go files:
  * [foundation/errors.go](foundation/errors.go)

## Gov

|Error Name|Codespace|Code|Description|
|:-|:-|:-|:-|
|ErrUnknownProposal|gov|2|unknown proposal|
|ErrInactiveProposal|gov|3|inactive proposal|
|ErrAlreadyActiveProposal|gov|4|proposal already active|
|ErrInvalidProposalContent|gov|5|invalid proposal content|
|ErrInvalidProposalType|gov|6|invalid proposal type|
|ErrInvalidVote|gov|7|invalid vote option|
|ErrInvalidGenesis|gov|8|invalid genesis state|
|ErrNoProposalHandlerExists|gov|9|no handler exists for proposal type|

>You can also find detailed information in the following Errors.go files:
  * [gov/types/errors.go](gov/types/errors.go)

## Params

|Error Name|Codespace|Code|Description|
|:-|:-|:-|:-|
|ErrUnknownSubspace|params|2|unknown subspace|
|ErrSettingParameter|params|3|failed to set parameter|
|ErrEmptyChanges|params|4|submitted parameter changes are empty|
|ErrEmptySubspace|params|5|parameter subspace is empty|
|ErrEmptyKey|params|6|parameter key is empty|
|ErrEmptyValue|params|7|parameter value is empty|

>You can also find detailed information in the following Errors.go files:
  * [params/types/proposal/errors.go](params/types/proposal/errors.go)

## Slashing

|Error Name|Codespace|Code|Description|
|:-|:-|:-|:-|
|ErrNoValidatorForAddress|slashing|2|address is not associated with any known validator|
|ErrBadValidatorAddr|slashing|3|validator does not exist for that address|
|ErrValidatorJailed|slashing|4|validator still jailed; cannot be unjailed|
|ErrValidatorNotJailed|slashing|5|validator not jailed; cannot be unjailed|
|ErrMissingSelfDelegation|slashing|6|validator has no self-delegation; cannot be unjailed|
|ErrSelfDelegationTooLowToUnjail|slashing|7|validator's self delegation less than minimum; cannot be unjailed|
|ErrNoSigningInfoFound|slashing|8|no validator signing info found|

>You can also find detailed information in the following Errors.go files:
  * [slashing/types/errors.go](slashing/types/errors.go)

## Staking

|Error Name|Codespace|Code|Description|
|:-|:-|:-|:-|
|ErrEmptyValidatorAddr|staking|2|empty validator address|
|ErrNoValidatorFound|staking|3|validator does not exist|
|ErrValidatorOwnerExists|staking|4|validator already exist for this operator address; must use new validator operator address|
|ErrValidatorPubKeyExists|staking|5|validator already exist for this pubkey; must use new validator pubkey|
|ErrValidatorPubKeyTypeNotSupported|staking|6|validator pubkey type is not supported|
|ErrValidatorJailed|staking|7|validator for this address is currently jailed|
|ErrBadRemoveValidator|staking|8|failed to remove validator|
|ErrCommissionNegative|staking|9|commission must be positive|
|ErrCommissionHuge|staking|10|commission cannot be more than 100%|
|ErrCommissionGTMaxRate|staking|11|commission cannot be more than the max rate|
|ErrCommissionUpdateTime|staking|12|commission cannot be changed more than once in 24h|
|ErrCommissionChangeRateNegative|staking|13|commission change rate must be positive|
|ErrCommissionChangeRateGTMaxRate|staking|14|commission change rate cannot be more than the max rate|
|ErrCommissionGTMaxChangeRate|staking|15|commission cannot be changed more than max change rate|
|ErrSelfDelegationBelowMinimum|staking|16|validator's self delegation must be greater than their minimum self delegation|
|ErrMinSelfDelegationDecreased|staking|17|minimum self delegation cannot be decrease|
|ErrEmptyDelegatorAddr|staking|18|empty delegator address|
|ErrNoDelegation|staking|19|no delegation for (address, validator) tuple|
|ErrBadDelegatorAddr|staking|20|delegator does not exist with address|
|ErrNoDelegatorForAddress|staking|21|delegator does not contain delegation|
|ErrInsufficientShares|staking|22|insufficient delegation shares|
|ErrDelegationValidatorEmpty|staking|23|cannot delegate to an empty validator|
|ErrNotEnoughDelegationShares|staking|24|not enough delegation shares|
|ErrNotMature|staking|25|entry not mature|
|ErrNoUnbondingDelegation|staking|26|no unbonding delegation found|
|ErrMaxUnbondingDelegationEntries|staking|27|too many unbonding delegation entries for (delegator, validator) tuple|
|ErrNoRedelegation|staking|28|no redelegation found|
|ErrSelfRedelegation|staking|29|cannot redelegate to the same validator|
|ErrTinyRedelegationAmount|staking|30|too few tokens to redelegate (truncates to zero tokens)|
|ErrBadRedelegationDst|staking|31|redelegation destination validator not found|
|ErrTransitiveRedelegation|staking|32|redelegation to this validator already in progress; first redelegation to this validator must complete before next redelegation|
|ErrMaxRedelegationEntries|staking|33|too many redelegation entries for (delegator, src-validator, dst-validator) tuple|
|ErrDelegatorShareExRateInvalid|staking|34|cannot delegate to validators with invalid (zero) ex-rate|
|ErrBothShareMsgsGiven|staking|35|both shares amount and shares percent provided|
|ErrNeitherShareMsgsGiven|staking|36|neither shares amount nor shares percent provided|
|ErrInvalidHistoricalInfo|staking|37|invalid historical info|
|ErrNoHistoricalInfo|staking|38|no historical info found|
|ErrEmptyValidatorPubKey|staking|39|empty validator public key|

>You can also find detailed information in the following Errors.go files:
  * [staking/types/errors.go](staking/types/errors.go)

## Token

|Error Name|Codespace|Code|Description|
|:-|:-|:-|:-|
|ErrInvalidContractID|contract|2|invalid contractID|
|ErrContractNotExist|contract|3|contract does not exist|
|ErrTokenNotExist|token|2|token does not exist|
|ErrTokenNotMintable|token|3|token is not mintable|
|ErrInvalidTokenName|token|4|token name should not be empty|
|ErrInvalidTokenDecimals|token|5|token decimals should be within the range in 0 ~ 18|
|ErrInvalidAmount|token|6|invalid token amount|
|ErrInvalidImageURILength|token|7|invalid token uri length|
|ErrInvalidNameLength|token|8|invalid name length|
|ErrInvalidTokenSymbol|token|9|invalid token symbol|
|ErrTokenNoPermission|token|10|account does not have the permission|
|ErrAccountExist|token|11|account already exists|
|ErrAccountNotExist|token|12|account does not exists|
|ErrInsufficientBalance|token|13|insufficient balance|
|ErrSupplyExist|token|14|supply for token already exists|
|ErrInsufficientSupply|token|15|insufficient supply|
|ErrInvalidChangesFieldCount|token|16|invalid count of field changes|
|ErrEmptyChanges|token|17|changes is empty|
|ErrInvalidChangesField|token|18|invalid field of changes|
|ErrDuplicateChangesField|token|19|invalid field of changes|
|ErrInvalidMetaLength|token|20|invalid meta length|
|ErrSupplyOverflow|token|21|supply for token reached maximum|
|ErrApproverProxySame|token|22|approver is same with proxy|
|ErrTokenNotApproved|token|23|proxy is not approved on the token|
|ErrTokenAlreadyApproved|token|24|proxy is already approved on the token|

>You can also find detailed information in the following Errors.go files:
  * [token/class/errors.go](token/class/errors.go)
  * [token/errors.go](token/errors.go)
