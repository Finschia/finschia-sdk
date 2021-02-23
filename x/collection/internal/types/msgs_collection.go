package types

import (
	"unicode/utf8"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
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

func (msg MsgCreateCollection) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner address cannot be empty")
	}
	if !ValidateName(msg.Name) {
		return sdkerrors.Wrapf(ErrInvalidNameLength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", msg.Name, MaxTokenNameLength, utf8.RuneCountInString(msg.Name))
	}
	if !ValidateMeta(msg.Meta) {
		return sdkerrors.Wrapf(ErrInvalidMetaLength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", msg.Meta, MaxTokenMetaLength, utf8.RuneCountInString(msg.Meta))
	}
	if !ValidateBaseImgURI(msg.BaseImgURI) {
		return sdkerrors.Wrapf(ErrInvalidBaseImgURILength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", msg.BaseImgURI, MaxBaseImgURILength, utf8.RuneCountInString(msg.BaseImgURI))
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
