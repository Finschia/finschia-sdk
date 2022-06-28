package collection_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	codectypes "github.com/line/lbm-sdk/codec/types"
	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
)

func TestValidateGenesis(t *testing.T) {
	addr := sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
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
					BaseImgUri: "file:///tibetian_fox.png",
				}},
			},
			false,
		},
		"contract of invalid name": {
			&collection.GenesisState{
				Contracts: []collection.Contract{{
					ContractId: "deadbeef",
					Name: string(make([]rune, 21)),
					Meta: "Tibetian Fox",
					BaseImgUri: "file:///tibetian_fox.png",
				}},
			},
			false,
		},
		"contract of invalid base img uri": {
			&collection.GenesisState{
				Contracts: []collection.Contract{{
					ContractId: "deadbeef",
					Name: "tibetian fox",
					BaseImgUri: string(make([]rune, 1001)),
					Meta: "Tibetian Fox",
				}},
			},
			false,
		},
		"contract of invalid meta": {
			&collection.GenesisState{
				Contracts: []collection.Contract{{
					ContractId: "deadbeef",
					Name: "tibetian fox",
					BaseImgUri: "file:///tibetian_fox.png",
					Meta: string(make([]rune, 1001)),
				}},
			},
			false,
		},
		"contract classes of invalid contract id": {
			&collection.GenesisState{
				Classes: []collection.ContractClasses{{
					Classes: []codectypes.Any{
						*collection.TokenClassToAny(&collection.NFTClass{
							Id: "deadbeef",
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
		"contract balances of invalid contract id": {
			&collection.GenesisState{
				Balances: []collection.ContractBalances{{
					Balances: []collection.Balance{{
						Address: addr.String(),
						Amount:  collection.NewCoins(collection.NewCoin("deadbeef" + fmt.Sprintf("%08x", 0), sdk.OneInt())),
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
						Amount:  collection.NewCoins(collection.NewCoin("deadbeef" + fmt.Sprintf("%08x", 0), sdk.OneInt())),
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
		"contract parents of invalid contract id": {
			&collection.GenesisState{
				Parents: []collection.ContractTokenRelations{{
					Relations: []collection.TokenRelation{{
						Self: "deadbeef",
						Other: "fee1dead",
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
						Other: "fee1dead" + fmt.Sprintf("%08x", 1),
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
						Self: "deadbeef" + fmt.Sprintf("%08x", 1),
					}},
				}},
			},
			false,
		},
		"contract authorizations of invalid contract id": {
			&collection.GenesisState{
				Authorizations: []collection.ContractAuthorizations{{
					Authorizations: []collection.Authorization{{
						Holder: addr.String(),
						Operator:    addr.String(),
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
						Operator:    addr.String(),
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
						Grantee: addr.String(),
						Permission: collection.Permission_Mint.String(),
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
						Permission: collection.Permission_Mint.String(),
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
