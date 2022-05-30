package collection

import (
	"fmt"

	proto "github.com/gogo/protobuf/proto"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

func DefaultNextClassIDs(contractID string) NextClassIDs {
	return NextClassIDs{
		ContractId:  contractID,
		Fungible:    sdk.NewUint(0),
		NonFungible: sdk.NewUint(1 << 28), // "10000000"
	}
}

type TokenClass interface {
	proto.Message

	GetContractId() string

	GetId() string
	SetId(ids *NextClassIDs)

	ValidateBasic() error
}

// FTClass
var _ TokenClass = (*FTClass)(nil)

//nolint
func (c *FTClass) SetId(ids *NextClassIDs) {
	id := ids.Fungible
	ids.Fungible = id.Incr()
	c.Id = fmt.Sprintf("%08x", id.Uint64())
}

func (c FTClass) ValidateBasic() error {
	if len(c.Id) != 8 || c.Id[0] != "0"[0] {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid id: %s", c.Id)
	}

	return nil
}

// NFTClass
var _ TokenClass = (*NFTClass)(nil)

//nolint
func (c *NFTClass) SetId(ids *NextClassIDs) {
	id := ids.NonFungible
	ids.NonFungible = id.Incr()
	c.Id = fmt.Sprintf("%08x", id.Uint64())
}

func (c NFTClass) ValidateBasic() error {
	if len(c.Id) != 8 || c.Id[0] == "0"[0] {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid id: %s", c.Id)
	}

	return nil
}
