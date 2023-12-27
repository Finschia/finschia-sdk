package foundation

import (
	"github.com/gogo/protobuf/proto"

	"cosmossdk.io/math"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:     DefaultParams(),
		Foundation: DefaultFoundation(),
	}
}

func DefaultFoundation() FoundationInfo {
	return *FoundationInfo{
		Version:     1,
		TotalWeight: math.LegacyZeroDec(),
	}.WithDecisionPolicy(DefaultDecisionPolicy())
}

func DefaultDecisionPolicy() DecisionPolicy {
	return &OutsourcingDecisionPolicy{
		Description: "using x/group",
	}
}

func DefaultAuthority() sdk.AccAddress {
	return authtypes.NewModuleAddress(ModuleName)
}

func DefaultParams() Params {
	return Params{
		FoundationTax: math.LegacyZeroDec(),
	}
}

func (data GenesisState) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if err := data.Foundation.UnpackInterfaces(unpacker); err != nil {
		return err
	}

	for _, p := range data.Proposals {
		if err := p.UnpackInterfaces(unpacker); err != nil {
			return err
		}
	}

	for _, ga := range data.Authorizations {
		if err := ga.UnpackInterfaces(unpacker); err != nil {
			return err
		}
	}

	return nil
}

func (i FoundationInfo) ValidateBasic() error {
	if i.Version == 0 {
		return sdkerrors.ErrInvalidVersion.Wrap("version must be > 0")
	}

	if i.TotalWeight.IsNil() || i.TotalWeight.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrap("total weight must be >= 0")
	}

	policy := i.GetDecisionPolicy()
	if policy == nil {
		return sdkerrors.ErrInvalidRequest.Wrap("must provide decision policy")
	}
	if err := policy.ValidateBasic(); err != nil {
		return err
	}

	// Is foundation outsourcing the proposal feature
	_, isOutsourcing := i.GetDecisionPolicy().(*OutsourcingDecisionPolicy)
	memberNotExists := i.TotalWeight.IsZero()
	if isOutsourcing && !memberNotExists {
		return sdkerrors.ErrInvalidRequest.Wrap("outsourcing policy not allows members")
	}
	if !isOutsourcing && memberNotExists {
		return sdkerrors.ErrInvalidRequest.Wrap("one member must exist at least")
	}

	return nil
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.ValidateBasic(); err != nil {
		return err
	}

	info := data.Foundation
	if err := info.ValidateBasic(); err != nil {
		return err
	}
	// Is x/foundation outsourcing the proposal feature
	isOutsourcing := info.TotalWeight.IsZero()

	if realWeight := math.LegacyNewDec(int64(len(data.Members))); !info.TotalWeight.Equal(realWeight) {
		return sdkerrors.ErrInvalidRequest.Wrapf("total weight not match, %s != %s", info.TotalWeight, realWeight)
	}
	members := Members{Members: data.Members}
	if err := members.ValidateBasic(); err != nil {
		return err
	}

	if isOutsourcing && len(data.Proposals) != 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("outsourcing policy not allows proposals")
	}
	proposalIDs := map[uint64]bool{}
	for _, proposal := range data.Proposals {
		id := proposal.Id
		if id > data.PreviousProposalId {
			return sdkerrors.ErrInvalidRequest.Wrapf("proposal %d has not yet been submitted", id)
		}
		if proposalIDs[id] {
			return sdkerrors.ErrInvalidRequest.Wrapf("duplicated proposal id of %d", id)
		}
		proposalIDs[id] = true

		if err := proposal.ValidateBasic(); err != nil {
			return err
		}

		if proposal.FoundationVersion > info.Version {
			return sdkerrors.ErrInvalidRequest.Wrapf("invalid foundation version of proposal %d", id)
		}
	}

	for _, vote := range data.Votes {
		if !proposalIDs[vote.ProposalId] {
			return sdkerrors.ErrInvalidRequest.Wrapf("vote for a proposal which does not exist: id %d", vote.ProposalId)
		}

		if _, err := sdk.AccAddressFromBech32(vote.Voter); err != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid voter address: %s", vote.Voter)
		}

		if err := validateVoteOption(vote.Option); err != nil {
			return err
		}
	}

	seenURLs := map[string]bool{}
	for _, censorship := range data.Censorships {
		if err := censorship.ValidateBasic(); err != nil {
			return err
		}
		if censorship.Authority == CensorshipAuthorityUnspecified {
			return sdkerrors.ErrInvalidRequest.Wrap("authority unspecified")
		}

		url := censorship.MsgTypeUrl
		if seenURLs[url] {
			return sdkerrors.ErrInvalidRequest.Wrapf("duplicate censorship over %s", url)
		}
		seenURLs[url] = true
	}

	for _, ga := range data.Authorizations {
		auth := ga.GetAuthorization()
		if auth == nil {
			return sdkerrors.ErrInvalidType.Wrap("invalid authorization")
		}

		url := auth.MsgTypeURL()
		if !seenURLs[url] {
			return sdkerrors.ErrInvalidRequest.Wrapf("no censorship over %s", url)
		}

		if _, err := sdk.AccAddressFromBech32(ga.Grantee); err != nil {
			return err
		}
	}

	if err := data.Pool.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

func (g GrantAuthorization) GetAuthorization() Authorization {
	if g.Authorization == nil {
		return nil
	}

	a, ok := g.Authorization.GetCachedValue().(Authorization)
	if !ok {
		return nil
	}
	return a
}

func (g *GrantAuthorization) SetAuthorization(a Authorization) error {
	msg, ok := a.(proto.Message)
	if !ok {
		return sdkerrors.ErrInvalidType.Wrapf("can't proto marshal %T", msg)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return err
	}
	g.Authorization = any

	return nil
}

// for the tests
func (g GrantAuthorization) WithAuthorization(authorization Authorization) *GrantAuthorization {
	grant := g
	if err := grant.SetAuthorization(authorization); err != nil {
		return nil
	}
	return &grant
}

func (g GrantAuthorization) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var authorization Authorization
	return unpacker.UnpackAny(g.Authorization, &authorization)
}
