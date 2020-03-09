package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = (*MsgCreateCollection)(nil)

type MsgCreateCollection struct {
	Owner      sdk.AccAddress `json:"owner"`
	Name       string         `json:"name"`
	Meta       string         `json:"meta"`
	BaseImgURI string         `json:"base_img_uri"`
}

func NewMsgCreateCollection(owner sdk.AccAddress, name, meta, baseImgURI string) MsgCreateCollection {
	return MsgCreateCollection{
		Owner:      owner,
		Name:       name,
		Meta:       meta,
		BaseImgURI: baseImgURI,
	}
}

func (msg MsgCreateCollection) ValidateBasic() sdk.Error {
	if len(msg.Name) == 0 {
		return ErrInvalidTokenName(DefaultCodespace, msg.Name)
	}
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("owner address cannot be empty")
	}
	if !ValidateBaseImgURI(msg.BaseImgURI) {
		return ErrInvalidBaseImgURILength(DefaultCodespace, msg.BaseImgURI)
	}
	return nil
}

func (MsgCreateCollection) Route() string { return RouterKey }
func (MsgCreateCollection) Type() string  { return "create_collection" }
func (msg MsgCreateCollection) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgCreateCollection) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
