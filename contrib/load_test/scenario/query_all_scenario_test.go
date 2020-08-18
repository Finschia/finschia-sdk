package scenario

import (
	"fmt"
	"testing"

	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	"github.com/stretchr/testify/require"
)

func TestQueryAllScenario_GenerateStateSettingMsgs(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	scenario, hdWallet, masterWallet := GivenTestEnvironments(t, server.URL, types.QueryAll, nil, nil)
	queryAllScenario, ok := scenario.(*QueryAllScenario)
	require.True(t, ok)

	msgs, params, err := queryAllScenario.GenerateStateSettingMsgs(masterWallet, hdWallet, []string{})
	require.NoError(t, err)

	require.Len(t, msgs, tests.TestTPS*tests.TestDuration*(8+13*tests.TestMsgsPerTxLoadTest))
	require.Equal(t, "send", msgs[tests.TestTPS*tests.TestDuration-1].Type())
	require.Equal(t, "grant_perm", msgs[tests.TestTPS*tests.TestDuration].Type())
	require.Equal(t, "grant_perm", msgs[8*tests.TestTPS*tests.TestDuration-1].Type())
	require.Equal(t, "mint_nft", msgs[8*tests.TestTPS*tests.TestDuration].Type())
	require.Equal(t, "9be17165", params["token_contract_id"])
	require.Equal(t, "678c146a", params["collection_contract_id"])
	require.Equal(t, "0000000100000000", params["ft_token_id"])
	require.Equal(t, "10000001", params["nft_token_type"])
	require.Equal(t, "16EFE7CF722157A57E03E947C6171B24A7FC3731E1A24FAE0D9168F80845407F", params["tx_hash"])
}

func TestQueryAllScenario_GenerateTarget(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	scenario, _, keyWallet := GivenTestEnvironments(t, server.URL, types.QueryAll,
		map[string]string{
			"token_contract_id":      tests.TokenContractID,
			"nft_token_type":         tests.NFTTokenType,
			"ft_token_id":            tests.FTTokenID,
			"tx_hash":                tests.TxHash,
			"collection_contract_id": tests.CollectionContractID,
			"num_nft_per_user":       fmt.Sprintf("%d", 13*tests.TestMsgsPerTxLoadTest),
		}, nil)
	queryAllScenario, ok := scenario.(*QueryAllScenario)
	require.True(t, ok)

	// When
	targets, numTargets, err := queryAllScenario.GenerateTarget(keyWallet, 0)
	require.NoError(t, err)

	// Then
	require.Equal(t, 97, numTargets)
	nftTokenID := fmt.Sprintf("%s%08x", tests.NFTTokenType, 1)
	expectedQueries := []string{
		"/supply/total",
		fmt.Sprintf("/token/%s/supply", tests.TokenContractID),
		fmt.Sprintf("/token/%s/token", tests.TokenContractID),
		fmt.Sprintf("/token/%s/accounts/%s/balance", tests.TokenContractID, tests.Address),
		fmt.Sprintf("/coin/balances/%s", tests.Address),
		"/genesis/app_state/accounts",
		fmt.Sprintf("/token/%s/supply", tests.TokenContractID),
		fmt.Sprintf("/collection/%s/fts/%s/supply", tests.CollectionContractID, tests.FTTokenID),
		fmt.Sprintf("/collection/%s/tokens", tests.CollectionContractID),
		"/coin/tcony",
		fmt.Sprintf("/token/%s/supply", tests.TokenContractID),
		fmt.Sprintf("/token/%s/token", tests.TokenContractID),
		fmt.Sprintf("/token/%s/accounts/%s/balance", tests.TokenContractID, tests.Address),
		fmt.Sprintf("/collection/%s/tokentypes", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/collection", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/tokentypes/%s/count", tests.CollectionContractID, tests.NFTTokenType),
		fmt.Sprintf("/collection/%s/fts/%s/mint", tests.CollectionContractID, tests.FTTokenID),
		fmt.Sprintf("/collection/%s/fts/%s/burn", tests.CollectionContractID, tests.FTTokenID),
		fmt.Sprintf("/collection/%s/fts/%s/supply", tests.CollectionContractID, tests.FTTokenID),
		fmt.Sprintf("/coin/balances/%s", tests.Address),
		"/genesis/app_state/accounts",
		fmt.Sprintf("/token/%s/supply", tests.TokenContractID),
		fmt.Sprintf("/collection/%s/fts/%s/supply", tests.CollectionContractID, tests.FTTokenID),
		fmt.Sprintf("/collection/%s/tokens", tests.CollectionContractID),
		"/coin/tcony",
		fmt.Sprintf("/token/%s/supply", tests.TokenContractID),
		fmt.Sprintf("/token/%s/token", tests.TokenContractID),
		fmt.Sprintf("/token/%s/accounts/%s/balance", tests.TokenContractID, tests.Address),
		fmt.Sprintf("/collection/%s/tokentypes", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/collection", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/tokentypes/%s/count", tests.CollectionContractID, tests.NFTTokenType),
		fmt.Sprintf("/collection/%s/fts/%s/mint", tests.CollectionContractID, tests.FTTokenID),
		fmt.Sprintf("/collection/%s/fts/%s/burn", tests.CollectionContractID, tests.FTTokenID),
		fmt.Sprintf("/collection/%s/fts/%s/supply", tests.CollectionContractID, tests.FTTokenID),
		fmt.Sprintf("/coin/balances/%s", tests.Address),
		"/genesis/app_state/accounts",
		fmt.Sprintf("/token/%s/supply", tests.TokenContractID),
		fmt.Sprintf("/collection/%s/fts/%s/supply", tests.CollectionContractID, tests.FTTokenID),
		fmt.Sprintf("/collection/%s/tokens", tests.CollectionContractID),
		"/coin/tcony",
		fmt.Sprintf("/token/%s/supply", tests.TokenContractID),
		fmt.Sprintf("/token/%s/token", tests.TokenContractID),
		fmt.Sprintf("/token/%s/accounts/%s/balance", tests.TokenContractID, tests.Address),
		fmt.Sprintf("/collection/%s/tokentypes", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/collection", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/tokentypes/%s/count", tests.CollectionContractID, tests.NFTTokenType),
		fmt.Sprintf("/collection/%s/fts/%s/mint", tests.CollectionContractID, tests.FTTokenID),
		fmt.Sprintf("/collection/%s/fts/%s/burn", tests.CollectionContractID, tests.FTTokenID),
		fmt.Sprintf("/collection/%s/fts/%s/supply", tests.CollectionContractID, tests.FTTokenID),
		fmt.Sprintf("/txs/%s", tests.TxHash),
		"/unconfirmed_txs",
		fmt.Sprintf("/collection/%s/nfts/%s/parent", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/collection/%s/nfts/%s/root", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/collection/%s/nfts/%s/children", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/txs/%s", tests.TxHash),
		"/unconfirmed_txs",
		fmt.Sprintf("/collection/%s/nfts/%s/parent", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/collection/%s/nfts/%s/root", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/collection/%s/nfts/%s/children", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/txs/%s", tests.TxHash),
		"/unconfirmed_txs",
		fmt.Sprintf("/collection/%s/nfts/%s/parent", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/collection/%s/nfts/%s/root", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/collection/%s/nfts/%s/children", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/txs/%s", tests.TxHash),
		"/unconfirmed_txs",
		fmt.Sprintf("/collection/%s/nfts/%s/parent", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/collection/%s/nfts/%s/root", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/collection/%s/nfts/%s/children", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/collection/%s/tokentypes", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/collection", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/tokentypes/%s/count", tests.CollectionContractID, tests.NFTTokenType),
		fmt.Sprintf("/collection/%s/nfts/%s", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/collection/%s/tokentypes", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/collection", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/tokentypes/%s/count", tests.CollectionContractID, tests.NFTTokenType),
		fmt.Sprintf("/collection/%s/nfts/%s", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/collection/%s/tokentypes", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/collection", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/tokentypes/%s/count", tests.CollectionContractID, tests.NFTTokenType),
		fmt.Sprintf("/collection/%s/nfts/%s", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/collection/%s/tokentypes", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/collection", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/tokentypes/%s/count", tests.CollectionContractID, tests.NFTTokenType),
		fmt.Sprintf("/collection/%s/nfts/%s", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/collection/%s/tokentypes", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/collection", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/tokentypes/%s/count", tests.CollectionContractID, tests.NFTTokenType),
		fmt.Sprintf("/collection/%s/nfts/%s", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/collection/%s/tokentypes", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/collection", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/tokentypes/%s/count", tests.CollectionContractID, tests.NFTTokenType),
		fmt.Sprintf("/collection/%s/nfts/%s", tests.CollectionContractID, nftTokenID),
		fmt.Sprintf("/collection/%s/tokentypes", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/collection", tests.CollectionContractID),
		fmt.Sprintf("/collection/%s/tokentypes/%s/count", tests.CollectionContractID, tests.NFTTokenType),
		fmt.Sprintf("/collection/%s/nfts/%s", tests.CollectionContractID, nftTokenID),
	}
	for i := 0; i < numTargets; i++ {
		require.Equal(t, "GET", (*targets)[i].Method)
		require.Equal(t, fmt.Sprintf("%s%s", server.URL, expectedQueries[i]), (*targets)[i].URL)
	}
}
