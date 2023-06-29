package types

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

var (
	ErrInvalidNameLength              = sdkerrors.Register(ModuleName, 1, "invalid name length")
	ErrSequencerExists                = sdkerrors.Register(ModuleName, 2, "sequencer already exist for this address")
	ErrNotExistRollupName             = sdkerrors.Register(ModuleName, 3, "rollup not exist for this name")
	SequencersByRollupExists          = sdkerrors.Register(ModuleName, 4, "sequencer already exist for this rollup")
	ErrPermissionedAddressesDuplicate = sdkerrors.Register(ModuleName, 5, "permissioned address has duplicates")
	ErrInvalidPermissionedAddress     = sdkerrors.Register(ModuleName, 6, "invalid permissioned address")
	ErrSequencerNotPermissioned       = sdkerrors.Register(ModuleName, 7, "sequencer is not permissioned for serving the rollup")
	ErrExistRollupName                = sdkerrors.Register(ModuleName, 8, "rollup already exist for this name")
)
