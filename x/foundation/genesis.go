package foundation

import (
	"github.com/gogo/protobuf/proto"
	codectypes "github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/authz"
)

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: &Params{
			Enabled: false,
		},
	}
}

func (data GenesisState) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, ga := range data.Authorizations {
		err := ga.UnpackInterfaces(unpacker)
		if err != nil {
			return err
		}
	}
	return nil
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if data.Params != nil {
		// validator auths are redundant where foundation is off
		if !data.Params.Enabled && len(data.ValidatorAuths) != 0 {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "redundant validator auths for disabled foundation")
		}

		if data.Params.FoundationTax.IsNegative() ||
			data.Params.FoundationTax.GT(sdk.OneDec()) {
			return sdkerrors.ErrInvalidRequest.Wrap("foundation tax must be >= 0 and <= 1")
		}
	}

	for _, auth := range data.ValidatorAuths {
		if err := sdk.ValidateValAddress(auth.OperatorAddress); err != nil {
			return err
		}
	}

	if info := data.Foundation; info != nil {
		if operator := info.Operator; len(operator) != 0 {
			if err := sdk.ValidateAccAddress(info.Operator); err != nil {
				return err
			}
		}

		if info.Version == 0 {
			return sdkerrors.ErrInvalidVersion.Wrap("version must be > 0")
		}

		if info.GetDecisionPolicy() != nil {
			if err := info.GetDecisionPolicy().ValidateBasic(); err != nil {
				return err
			}
		}

	}

	if err := validateMembers(data.Members); err != nil {
		return err
	}

	proposalIDs := map[uint64]bool{}
	for _, proposal := range data.Proposals {
		id := proposal.Id
		if proposalIDs[id] {
			return sdkerrors.ErrInvalidRequest.Wrapf("duplicated id: %d", id)
		}
		proposalIDs[id] = true

		if err := validateProposers(proposal.Proposers); err != nil {
			return err
		}

		if proposal.FoundationVersion == 0 {
			return sdkerrors.ErrInvalidVersion.Wrap("foundation version must be > 0")
		}

		if err := validateMsgs(proposal.GetMsgs()); err != nil {
			return err
		}
	}

	for _, vote := range data.Votes {
		if !proposalIDs[vote.ProposalId] {
			return sdkerrors.ErrInvalidRequest.Wrapf("vote for a proposal which does not exist: id %d", vote.ProposalId)
		}

		if err := sdk.ValidateAccAddress(vote.Voter); err != nil {
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

		if err := sdk.ValidateAccAddress(ga.Grantee); err != nil {
			return err
		}
	}

	return nil
}

func (g GrantAuthorization) GetAuthorization() authz.Authorization {
	if g.Authorization == nil {
		return nil
	}

	a, ok := g.Authorization.GetCachedValue().(authz.Authorization)
	if !ok {
		return nil
	}
	return a
}

func (g *GrantAuthorization) SetAuthorization(a authz.Authorization) error {
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
func (g GrantAuthorization) WithAuthorization(authorization authz.Authorization) *GrantAuthorization {
	grant := g
	if err := grant.SetAuthorization(authorization); err != nil {
		return nil
	}
	return &grant
}

func (g GrantAuthorization) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var authorization authz.Authorization
	return unpacker.UnpackAny(g.Authorization, &authorization)
}
