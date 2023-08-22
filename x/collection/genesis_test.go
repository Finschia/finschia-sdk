package collection_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	codectypes "github.com/Finschia/finschia-rdk/codec/types"
	"github.com/Finschia/finschia-rdk/crypto/keys/secp256k1"
	sdk "github.com/Finschia/finschia-rdk/types"
	"github.com/Finschia/finschia-rdk/x/collection"
)

func TestValidateGenesis(t *testing.T) {
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	testCases := map[string]struct {
		gs    *collection.GenesisState
		valid bool
	}{
		"default genesis": {
			collection.DefaultGenesisState(),
			true,
		},
		"contract of invalid contract id": {
			&collection.GenesisState{
				Contracts: []collection.Contract{{
					Name: "tibetian fox",
					Meta: "Tibetian Fox",
					Uri:  "file:///tibetian_fox.png",
				}},
			},
			false,
		},
		"contract of invalid name": {
			&collection.GenesisState{
				Contracts: []collection.Contract{{
					Id:   "deadbeef",
					Name: string(make([]rune, 21)),
					Meta: "Tibetian Fox",
					Uri:  "file:///tibetian_fox.png",
				}},
			},
			false,
		},
		"contract of invalid base img uri": {
			&collection.GenesisState{
				Contracts: []collection.Contract{{
					Id:   "deadbeef",
					Name: "tibetian fox",
					Uri:  string(make([]rune, 1001)),
					Meta: "Tibetian Fox",
				}},
			},
			false,
		},
		"contract of invalid meta": {
			&collection.GenesisState{
				Contracts: []collection.Contract{{
					Id:   "deadbeef",
					Name: "tibetian fox",
					Uri:  "file:///tibetian_fox.png",
					Meta: string(make([]rune, 1001)),
				}},
			},
			false,
		},
		"next class ids of invalid contract id": {
			&collection.GenesisState{
				NextClassIds: []collection.NextClassIDs{{
					Fungible:    sdk.ZeroUint(),
					NonFungible: sdk.OneUint(),
				}},
			},
			false,
		},
		"contract classes of invalid contract id": {
			&collection.GenesisState{
				Classes: []collection.ContractClasses{{
					Classes: []codectypes.Any{
						*collection.TokenClassToAny(&collection.NFTClass{
							Id:   "deadbeef",
							Name: "tibetian fox",
							Meta: "Tibetian Fox",
						}),
					},
				}},
			},
			false,
		},
		"contract classes of empty classes": {
			&collection.GenesisState{
				Classes: []collection.ContractClasses{{
					ContractId: "deadbeef",
				}},
			},
			false,
		},
		"contract classes of invalid class": {
			&collection.GenesisState{
				Classes: []collection.ContractClasses{{
					ContractId: "deadbeef",
					Classes: []codectypes.Any{
						*collection.TokenClassToAny(&collection.NFTClass{
							Name: "tibetian fox",
							Meta: "Tibetian Fox",
						}),
					},
				}},
			},
			false,
		},
		"contract next token ids of invalid contract id": {
			&collection.GenesisState{
				NextTokenIds: []collection.ContractNextTokenIDs{{
					TokenIds: []collection.NextTokenID{{
						ClassId: "deadbeef",
						Id:      sdk.ZeroUint(),
					}},
				}},
			},
			false,
		},
		"contract next token ids of empty classes": {
			&collection.GenesisState{
				NextTokenIds: []collection.ContractNextTokenIDs{{
					ContractId: "deadbeef",
				}},
			},
			false,
		},
		"contract next token ids of invalid class": {
			&collection.GenesisState{
				NextTokenIds: []collection.ContractNextTokenIDs{{
					ContractId: "deadbeef",
					TokenIds: []collection.NextTokenID{{
						Id: sdk.ZeroUint(),
					}},
				}},
			},
			false,
		},
		"contract balances of invalid contract id": {
			&collection.GenesisState{
				Balances: []collection.ContractBalances{{
					Balances: []collection.Balance{{
						Address: addr.String(),
						Amount:  collection.NewCoins(collection.NewFTCoin("00bab10c", sdk.OneInt())),
					}},
				}},
			},
			false,
		},
		"contract balances of empty balances": {
			&collection.GenesisState{
				Balances: []collection.ContractBalances{{
					ContractId: "deadbeef",
				}},
			},
			false,
		},
		"contract balances of invalid address": {
			&collection.GenesisState{
				Balances: []collection.ContractBalances{{
					ContractId: "deadbeef",
					Balances: []collection.Balance{{
						Amount: collection.NewCoins(collection.NewFTCoin("00bab10c", sdk.OneInt())),
					}},
				}},
			},
			false,
		},
		"contract balances of invalid amount": {
			&collection.GenesisState{
				Balances: []collection.ContractBalances{{
					ContractId: "deadbeef",
					Balances: []collection.Balance{{
						Address: addr.String(),
					}},
				}},
			},
			false,
		},
		"contract nfts of invalid contract id": {
			&collection.GenesisState{
				Nfts: []collection.ContractNFTs{{
					Nfts: []collection.NFT{{
						TokenId: collection.NewNFTID("deadbeef", 1),
						Name:    "tibetian fox",
						Meta:    "Tibetian Fox",
					}},
				}},
			},
			false,
		},
		"contract nfts of empty nfts": {
			&collection.GenesisState{
				Nfts: []collection.ContractNFTs{{
					ContractId: "deadbeef",
				}},
			},
			false,
		},
		"contract nfts of invalid class": {
			&collection.GenesisState{
				Nfts: []collection.ContractNFTs{{
					ContractId: "deadbeef",
					Nfts: []collection.NFT{{
						Name: "tibetian fox",
						Meta: "Tibetian Fox",
					}},
				}},
			},
			false,
		},
		"contract nfts of invalid name": {
			&collection.GenesisState{
				Nfts: []collection.ContractNFTs{{
					ContractId: "deadbeef",
					Nfts: []collection.NFT{{
						TokenId: collection.NewNFTID("deadbeef", 1),
						Name:    string(make([]rune, 21)),
						Meta:    "Tibetian Fox",
					}},
				}},
			},
			false,
		},
		"contract nfts of invalid meta": {
			&collection.GenesisState{
				Nfts: []collection.ContractNFTs{{
					ContractId: "deadbeef",
					Nfts: []collection.NFT{{
						TokenId: collection.NewNFTID("deadbeef", 1),
						Name:    "tibetian fox",
						Meta:    string(make([]rune, 1001)),
					}},
				}},
			},
			false,
		},
		"contract parents of invalid contract id": {
			&collection.GenesisState{
				Parents: []collection.ContractTokenRelations{{
					Relations: []collection.TokenRelation{{
						Self:  collection.NewNFTID("deadbeef", 1),
						Other: collection.NewNFTID("fee1dead", 1),
					}},
				}},
			},
			false,
		},
		"contract parents of empty relations": {
			&collection.GenesisState{
				Parents: []collection.ContractTokenRelations{{
					ContractId: "deadbeef",
				}},
			},
			false,
		},
		"contract parents of invalid token": {
			&collection.GenesisState{
				Parents: []collection.ContractTokenRelations{{
					ContractId: "deadbeef",
					Relations: []collection.TokenRelation{{
						Other: collection.NewNFTID("fee1dead", 1),
					}},
				}},
			},
			false,
		},
		"contract parents of invalid parent": {
			&collection.GenesisState{
				Parents: []collection.ContractTokenRelations{{
					ContractId: "deadbeef",
					Relations: []collection.TokenRelation{{
						Self: collection.NewNFTID("deadbeef", 1),
					}},
				}},
			},
			false,
		},
		"contract authorizations of invalid contract id": {
			&collection.GenesisState{
				Authorizations: []collection.ContractAuthorizations{{
					Authorizations: []collection.Authorization{{
						Holder:   addr.String(),
						Operator: addr.String(),
					}},
				}},
			},
			false,
		},
		"contract authorizations of empty authorizations": {
			&collection.GenesisState{
				Authorizations: []collection.ContractAuthorizations{{
					ContractId: "deadbeef",
				}},
			},
			false,
		},
		"contract authorizations of invalid holder": {
			&collection.GenesisState{
				Authorizations: []collection.ContractAuthorizations{{
					ContractId: "deadbeef",
					Authorizations: []collection.Authorization{{
						Operator: addr.String(),
					}},
				}},
			},
			false,
		},
		"contract authorizations of invalid operator": {
			&collection.GenesisState{
				Authorizations: []collection.ContractAuthorizations{{
					ContractId: "deadbeef",
					Authorizations: []collection.Authorization{{
						Holder: addr.String(),
					}},
				}},
			},
			false,
		},
		"contract grants of invalid contract id": {
			&collection.GenesisState{
				Grants: []collection.ContractGrants{{
					Grants: []collection.Grant{{
						Grantee:    addr.String(),
						Permission: collection.PermissionMint,
					}},
				}},
			},
			false,
		},
		"contract grants of empty grants": {
			&collection.GenesisState{
				Grants: []collection.ContractGrants{{
					ContractId: "deadbeef",
				}},
			},
			false,
		},
		"contract grants of invalid grantee": {
			&collection.GenesisState{
				Grants: []collection.ContractGrants{{
					ContractId: "deadbeef",
					Grants: []collection.Grant{{
						Permission: collection.PermissionMint,
					}},
				}},
			},
			false,
		},
		"contract grants of invalid permission": {
			&collection.GenesisState{
				Grants: []collection.ContractGrants{{
					ContractId: "deadbeef",
					Grants: []collection.Grant{{
						Grantee: addr.String(),
					}},
				}},
			},
			false,
		},
		"contract supplies of invalid contract id": {
			&collection.GenesisState{
				Supplies: []collection.ContractStatistics{{
					Statistics: []collection.ClassStatistics{{
						ClassId: "deadbeef",
						Amount:  sdk.OneInt(),
					}},
				}},
			},
			false,
		},
		"contract supplies of empty supplies": {
			&collection.GenesisState{
				Supplies: []collection.ContractStatistics{{
					ContractId: "deadbeef",
				}},
			},
			false,
		},
		"contract supplies of invalid class id": {
			&collection.GenesisState{
				Supplies: []collection.ContractStatistics{{
					ContractId: "deadbeef",
					Statistics: []collection.ClassStatistics{{
						Amount: sdk.OneInt(),
					}},
				}},
			},
			false,
		},
		"contract supplies of invalid operator": {
			&collection.GenesisState{
				Supplies: []collection.ContractStatistics{{
					ContractId: "deadbeef",
					Statistics: []collection.ClassStatistics{{
						ClassId: "deadbeef",
						Amount:  sdk.ZeroInt(),
					}},
				}},
			},
			false,
		},
		"contract burnts of invalid contract id": {
			&collection.GenesisState{
				Burnts: []collection.ContractStatistics{{
					Statistics: []collection.ClassStatistics{{
						ClassId: "deadbeef",
						Amount:  sdk.OneInt(),
					}},
				}},
			},
			false,
		},
		"contract burnts of empty burnts": {
			&collection.GenesisState{
				Burnts: []collection.ContractStatistics{{
					ContractId: "deadbeef",
				}},
			},
			false,
		},
		"contract burnts of invalid class id": {
			&collection.GenesisState{
				Burnts: []collection.ContractStatistics{{
					ContractId: "deadbeef",
					Statistics: []collection.ClassStatistics{{
						Amount: sdk.OneInt(),
					}},
				}},
			},
			false,
		},
		"contract burnts of invalid operator": {
			&collection.GenesisState{
				Burnts: []collection.ContractStatistics{{
					ContractId: "deadbeef",
					Statistics: []collection.ClassStatistics{{
						ClassId: "deadbeef",
						Amount:  sdk.ZeroInt(),
					}},
				}},
			},
			false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := collection.ValidateGenesis(*tc.gs)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
