package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
)

var _ sdk.Msg = (*MsgIssue)(nil)
var _ sdk.Msg = (*MsgIssueCollection)(nil)
var _ sdk.Msg = (*MsgMint)(nil)
var _ sdk.Msg = (*MsgBurn)(nil)
var _ sdk.Msg = (*MsgGrantPermission)(nil)
var _ sdk.Msg = (*MsgRevokePermission)(nil)

var _ json.Marshaler = (*MsgIssue)(nil)
var _ json.Unmarshaler = (*MsgIssue)(nil)
var _ json.Marshaler = (*MsgIssueCollection)(nil)
var _ json.Unmarshaler = (*MsgIssueCollection)(nil)

type MsgIssue struct {
	MsgIssueCommon
	Amount   sdk.Int `json:"amount"`
	Mintable bool    `json:"mintable"`
	Decimals sdk.Int `json:"decimals"`
}

func NewMsgIssue(name, symbol, tokenURI string, owner sdk.AccAddress, amount sdk.Int, decimal sdk.Int, mintable bool) MsgIssue {
	return MsgIssue{
		MsgIssueCommon: NewMsgIssueCommon(name, symbol, tokenURI, owner),
		Amount:         amount,
		Mintable:       mintable,
		Decimals:       decimal,
	}
}

// ATTENTION: to avoid embedded serialization implement this custom JSON marshaler
// see https://github.com/tendermint/go-amino/issues/269
func (msg MsgIssue) MarshalJSON() ([]byte, error) {
	type msgAlias MsgIssue
	return json.Marshal((msgAlias)(msg))
}

func (msg *MsgIssue) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgIssue
	return json.Unmarshal(data, msgAlias(msg))
}

func (msg MsgIssue) Route() string { return RouterKey }

func (msg MsgIssue) Type() string { return "issue_token" }

func (msg MsgIssue) ValidateBasic() sdk.Error {
	if err := msg.MsgIssueCommon.ValidateBasic(); err != nil {
		return err
	}

	if msg.Amount.Equal(sdk.NewInt(1)) && msg.Decimals.Equal(sdk.NewInt(0)) && !msg.Mintable {
		return ErrInvalidFTExist(DefaultCodespace)
	}

	if msg.Decimals.GT(sdk.NewInt(18)) || msg.Decimals.IsNegative() {
		return ErrTokenInvalidDecimals(DefaultCodespace)
	}

	coin := sdk.NewCoin(msg.Symbol, msg.Amount)
	if !coin.IsValid() {
		return sdk.ErrInvalidCoins(coin.String())
	}

	return nil
}

func (msg MsgIssue) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIssue) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

type MsgIssueCollection struct {
	MsgIssue
	MsgCollection
}

func NewMsgIssueCollection(name, symbol, tokenURI string, owner sdk.AccAddress, amount sdk.Int, decimal sdk.Int, mintable bool, tokenID string) MsgIssueCollection {
	return MsgIssueCollection{
		MsgIssue:      NewMsgIssue(name, symbol, tokenURI, owner, amount, decimal, mintable),
		MsgCollection: NewMsgCollection(tokenID),
	}
}

// ATTENTION: to avoid embedded serialization implement this custom JSON marshaler
// see https://github.com/tendermint/go-amino/issues/269
func (msg MsgIssueCollection) MarshalJSON() ([]byte, error) {
	type msgAlias MsgIssueCollection
	return json.Marshal((msgAlias)(msg))
}

func (msg *MsgIssueCollection) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgIssueCollection
	return json.Unmarshal(data, msgAlias(msg))
}

func (msg MsgIssueCollection) Type() string { return "issue_token_collection" }

func (msg MsgIssueCollection) ValidateBasic() sdk.Error {
	if err := msg.MsgIssue.ValidateBasic(); err != nil {
		return err
	}

	if err := msg.MsgCollection.ValidateBasic(); err != nil {
		return err
	}

	coin := sdk.NewCoin(types.SymbolCollectionToken(msg.Symbol, msg.TokenID), msg.Amount)
	if !coin.IsValid() {
		return sdk.ErrInvalidCoins(coin.String())
	}
	return nil
}

func NewMsgMint(to sdk.AccAddress, amount sdk.Coins) MsgMint {
	return MsgMint{
		To:     to,
		Amount: amount,
	}
}

type MsgMint struct {
	To     sdk.AccAddress `json:"to"`
	Amount sdk.Coins      `json:"amount"`
}

func (MsgMint) Route() string { return RouterKey }

func (MsgMint) Type() string { return "mint_token" }

func (msg MsgMint) ValidateBasic() sdk.Error {
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("amount is not valid")
	}
	if msg.To.Empty() {
		return sdk.ErrInvalidAddress("To address cannot be empty")
	}
	return nil
}

func (msg MsgMint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgMint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.To}
}

func NewMsgBurn(from sdk.AccAddress, amount sdk.Coins) MsgBurn {
	return MsgBurn{
		From:   from,
		Amount: amount,
	}
}

type MsgBurn struct {
	From   sdk.AccAddress `json:"from"`
	Amount sdk.Coins      `json:"amount"`
}

func (MsgBurn) Route() string { return RouterKey }

func (MsgBurn) Type() string { return "burn_token" }

func (msg MsgBurn) ValidateBasic() sdk.Error {
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("amount is not valid")
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From address cannot be empty")
	}
	return nil
}

func (msg MsgBurn) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func NewMsgGrantPermission(from, to sdk.AccAddress, perm Permission) MsgGrantPermission {
	return MsgGrantPermission{
		From:       from,
		To:         to,
		Permission: perm,
	}
}

type MsgGrantPermission struct {
	From       sdk.AccAddress `json:"from"`
	To         sdk.AccAddress `json:"to"`
	Permission Permission     `json:"permission"`
}

func (MsgGrantPermission) Route() string { return RouterKey }

func (MsgGrantPermission) Type() string { return "grant_permission" }

func (msg MsgGrantPermission) ValidateBasic() sdk.Error {
	if msg.From.Empty() || msg.To.Empty() {
		return sdk.ErrInvalidAddress("addresses cannot be empty")
	}

	if msg.From.Equals(msg.To) {
		return sdk.ErrInvalidAddress("from, to address can not be the same")
	}

	if len(msg.Permission.GetAction()) <= 0 || len(msg.Permission.GetResource()) <= 0 {
		return ErrTokenPermission(DefaultCodespace)
	}
	return nil
}

func (msg MsgGrantPermission) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgGrantPermission) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func NewMsgRevokePermission(from sdk.AccAddress, perm Permission) MsgRevokePermission {
	return MsgRevokePermission{
		From:       from,
		Permission: perm,
	}
}

type MsgRevokePermission struct {
	From       sdk.AccAddress `json:"from"`
	Permission Permission     `json:"permission"`
}

func (MsgRevokePermission) Route() string { return RouterKey }

func (MsgRevokePermission) Type() string { return "revoke_permission" }

func (msg MsgRevokePermission) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("addresses cannot be empty")
	}

	if len(msg.Permission.GetAction()) <= 0 || len(msg.Permission.GetResource()) <= 0 {
		return ErrTokenPermission(DefaultCodespace)
	}
	return nil
}

func (msg MsgRevokePermission) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRevokePermission) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

type MsgIssueCommon struct {
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	Owner    sdk.AccAddress `json:"owner"`
	TokenURI string         `json:"token_uri"`
}

func NewMsgIssueCommon(name, symbol, tokenURI string, owner sdk.AccAddress) MsgIssueCommon {
	return MsgIssueCommon{
		Name:     name,
		Symbol:   symbol,
		TokenURI: tokenURI,
		Owner:    owner,
	}
}
func (msg MsgIssueCommon) ValidateBasic() sdk.Error {
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace)
	}
	if len(msg.Name) == 0 {
		return ErrInvalidTokenName(DefaultCodespace)
	}
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("owner address cannot be empty")
	}
	return nil
}

type MsgCollection struct {
	TokenID string `json:"token_id"`
}

// ATTENTION: to avoid embedded serialization implement this custom JSON marshaler
// see https://github.com/tendermint/go-amino/issues/269
func (msg MsgCollection) MarshalJSON() ([]byte, error) {
	type msgAlias MsgCollection
	return json.Marshal((msgAlias)(msg))
}

func (msg *MsgCollection) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgCollection
	return json.Unmarshal(data, msgAlias(msg))
}

func (msg MsgCollection) ValidateBasic() sdk.Error {
	if err := types.ValidateSymbolTokenID(msg.TokenID); err != nil {
		return ErrInvalidTokenID(DefaultCodespace)
	}
	return nil
}

func NewMsgCollection(tokenID string) MsgCollection {
	return MsgCollection{
		TokenID: tokenID,
	}
}
