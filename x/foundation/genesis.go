package foundation

import (
	"github.com/gogo/protobuf/proto"
	codectypes "github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
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
		Operator:    DefaultOperator().String(),
		Version:     1,
		TotalWeight: sdk.ZeroDec(),
	}.WithDecisionPolicy(DefaultDecisionPolicy())
}

func DefaultOperator() sdk.AccAddress {
	return authtypes.NewModuleAddress(DefaultOperatorName)
}

func DefaultParams() Params {
	return Params{
		FoundationTax: sdk.ZeroDec(),
	}
}

func (data GenesisState) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if err := data.Foundation.UnpackInterfaces(unpacker); err != nil {
		return err
	}

	for _, ga := range data.Authorizations {
		if err := ga.UnpackInterfaces(unpacker); err != nil {
			return err
		}
	}
	return nil
}

func (i FoundationInfo) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(i.Operator); err != nil {
		return err
	}

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
	if realWeight := sdk.NewDecFromInt(sdk.NewInt(int64(len(data.Members)))); !info.TotalWeight.Equal(realWeight) {
		return sdkerrors.ErrInvalidRequest.Wrapf("total weight not match, %s != %s", info.TotalWeight, realWeight)
	}

	_, outsourcing := info.GetDecisionPolicy().(*OutsourcingDecisionPolicy)

	if outsourcing && len(data.Members) != 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("outsourcing policy not allows members")
	}
	members := Members{Members: data.Members}
	if err := members.ValidateBasic(); err != nil {
		return err
	}

	if outsourcing && len(data.Proposals) != 0 {
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

	for _, ga := range data.Authorizations {
		if ga.GetAuthorization() == nil {
			return sdkerrors.ErrInvalidType.Wrap("invalid authorization")
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
