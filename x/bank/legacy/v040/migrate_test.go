package v040_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lfb-sdk/client"
	"github.com/line/lfb-sdk/simapp"
	sdk "github.com/line/lfb-sdk/types"
	v038auth "github.com/line/lfb-sdk/x/auth/legacy/v038"
	v039auth "github.com/line/lfb-sdk/x/auth/legacy/v039"
	v036supply "github.com/line/lfb-sdk/x/bank/legacy/v036"
	v038bank "github.com/line/lfb-sdk/x/bank/legacy/v038"
	v040bank "github.com/line/lfb-sdk/x/bank/legacy/v040"
)

func TestMigrate(t *testing.T) {
	encodingConfig := simapp.MakeTestEncodingConfig()
	clientCtx := client.Context{}.
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithJSONMarshaler(encodingConfig.Marshaler)

	coins := sdk.NewCoins(sdk.NewInt64Coin("stake", 50))
	addr1, _ := sdk.AccAddressFromBech32("link17dgvcdx0v4mlxfrmfhua7685py3akv3lpnlpce")
	acc1 := v038auth.NewBaseAccount(addr1, coins, nil, 1, 0)

	addr2, _ := sdk.AccAddressFromBech32("link15lclrh8eqj233lmvxj4kcut2mua03t7u09ff00")
	vaac := v038auth.NewContinuousVestingAccountRaw(
		v038auth.NewBaseVestingAccount(
			v038auth.NewBaseAccount(addr2, coins, nil, 1, 0), coins, nil, nil, 3160620846,
		),
		1580309972,
	)

	supply := sdk.NewCoins(sdk.NewInt64Coin("stake", 1000))

	bankGenState := v038bank.GenesisState{
		SendEnabled: true,
	}
	authGenState := v039auth.GenesisState{
		Accounts: v038auth.GenesisAccounts{acc1, vaac},
	}
	supplyGenState := v036supply.GenesisState{
		Supply: supply,
	}

	migrated := v040bank.Migrate(bankGenState, authGenState, supplyGenState)
	expected := `{"params":{"send_enabled":[],"default_send_enabled":true},"balances":[{"address":"link17dgvcdx0v4mlxfrmfhua7685py3akv3lpnlpce","coins":[{"denom":"stake","amount":"50"}]},{"address":"link15lclrh8eqj233lmvxj4kcut2mua03t7u09ff00","coins":[{"denom":"stake","amount":"50"}]}],"supply":[{"denom":"stake","amount":"1000"}],"denom_metadata":[]}`

	bz, err := clientCtx.JSONMarshaler.MarshalJSON(migrated)
	require.NoError(t, err)
	require.Equal(t, expected, string(bz))
}
