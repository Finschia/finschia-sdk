package scenario

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/transaction"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	"github.com/line/link/x/coin"
	"github.com/line/link/x/collection"
	"github.com/line/link/x/token"
)

func GenerateRegisterAccountMsgs(masterAddress sdk.AccAddress, hdWallet *wallet.HDWallet,
	config types.Config) ([]sdk.Msg, error) {
	numMsgs := config.TPS * config.Duration
	coins := sdk.NewCoins(sdk.NewCoin(config.CoinName, sdk.NewInt(1000*int64(config.MsgsPerTxLoadTest))))
	msgs := make([]sdk.Msg, numMsgs)
	for i := 0; i < numMsgs; i++ {
		keyWallet, err := hdWallet.GetKeyWallet(0, uint32(i))
		if err != nil {
			return nil, err
		}
		to := keyWallet.Address()
		msgs[i] = coin.NewMsgSend(masterAddress, to, coins)
	}
	return msgs, nil
}

func GenerateGrantPermissionMsgs(masterAddress sdk.AccAddress, hdWallet *wallet.HDWallet,
	config types.Config, contractID string, module string, perms []string) ([]sdk.Msg, error) {
	numUsers := config.TPS * config.Duration
	msgs := make([]sdk.Msg, numUsers*len(perms))
	for i := 0; i < numUsers; i++ {
		keyWallet, err := hdWallet.GetKeyWallet(0, uint32(i))
		if err != nil {
			return nil, err
		}
		to := keyWallet.Address()
		for j := 0; j < len(perms); j++ {
			switch module {
			case "token":
				msgs[len(perms)*i+j] = token.NewMsgGrantPermission(masterAddress, contractID, to,
					token.Permission(perms[j]))
			case "collection":
				msgs[len(perms)*i+j] = collection.NewMsgGrantPermission(masterAddress, contractID, to,
					collection.Permission(perms[j]))
			default:
				return nil, types.InvalidModuleName{Name: module}
			}
		}
	}
	return msgs, nil
}

func GenerateMintFTMsgs(masterAddress sdk.AccAddress, hdWallet *wallet.HDWallet,
	config types.Config, contractID, ftTokenID string) ([]sdk.Msg, error) {
	numUsers := config.TPS * config.Duration
	msgs := make([]sdk.Msg, numUsers)
	for i := 0; i < numUsers; i++ {
		keyWallet, err := hdWallet.GetKeyWallet(0, uint32(i))
		if err != nil {
			return nil, err
		}
		to := keyWallet.Address()
		msgs[i] = collection.NewMsgMintFT(masterAddress, contractID, to, collection.NewCoin(ftTokenID, sdk.NewInt(10)))
	}
	return msgs, nil
}

func GenerateMintNFTMsgs(masterAddress sdk.AccAddress, hdWallet *wallet.HDWallet,
	config types.Config, contractID, tokenType string, numNFTPerUser int) ([]sdk.Msg, error) {
	numUsers := config.TPS * config.Duration
	msgs := make([]sdk.Msg, numUsers*numNFTPerUser)
	for i := 0; i < numUsers; i++ {
		keyWallet, err := hdWallet.GetKeyWallet(0, uint32(i))
		if err != nil {
			return nil, err
		}
		to := keyWallet.Address()
		for j := 0; j < numNFTPerUser; j++ {
			msgs[numNFTPerUser*i+j] = collection.NewMsgMintNFT(masterAddress, contractID, to,
				collection.NewMintNFTParam("name", "{}", tokenType))
		}
	}
	return msgs, nil
}

func IssueToken(linkService *service.LinkService, txBuilder transaction.TxBuilderWithoutKeybase,
	masterKeyWallet *wallet.KeyWallet) (string, string, error) {
	masterAddress := masterKeyWallet.Address()
	account, err := linkService.GetAccount(masterAddress.String())
	if err != nil {
		return "", "", err
	}
	res, err := linkService.BroadcastMsgs(
		txBuilder,
		account,
		masterKeyWallet,
		[]sdk.Msg{token.NewMsgIssue(masterAddress, masterAddress, "token", "TK", "{}", "imageuri",
			sdk.NewInt(1), sdk.NewInt(8), true)},
		service.BroadcastBlock,
	)
	if err != nil {
		return "", "", err
	}
	if res.Code != 0 {
		return "", res.TxHash, types.FailedTxError{Tx: res}
	}
	for _, attribute := range res.Logs[0].Events[0].Attributes {
		if attribute.Key == "contract_id" {
			return attribute.Value, res.TxHash, nil
		}
	}
	return "", res.TxHash, types.NoContractIDError{Tx: res}
}

func CreateCollection(linkService *service.LinkService, txBuilder transaction.TxBuilderWithoutKeybase,
	masterKeyWallet *wallet.KeyWallet) (string, error) {
	masterAddress := masterKeyWallet.Address()
	account, err := linkService.GetAccount(masterAddress.String())
	if err != nil {
		return "", err
	}
	res, err := linkService.BroadcastMsgs(
		txBuilder,
		account,
		masterKeyWallet,
		[]sdk.Msg{collection.NewMsgCreateCollection(masterAddress, "collection", "{}", "uri")},
		service.BroadcastBlock,
	)
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", types.FailedTxError{Tx: res}
	}
	for _, attribute := range res.Logs[0].Events[0].Attributes {
		if attribute.Key == "contract_id" {
			return attribute.Value, nil
		}
	}
	return "", types.NoContractIDError{Tx: res}
}

func IssueFT(linkService *service.LinkService, txBuilder transaction.TxBuilderWithoutKeybase,
	masterKeyWallet *wallet.KeyWallet, contractID string) (string, error) {
	masterAddress := masterKeyWallet.Address()
	account, err := linkService.GetAccount(masterAddress.String())
	if err != nil {
		return "", err
	}
	res, err := linkService.BroadcastMsgs(
		txBuilder,
		account,
		masterKeyWallet,
		[]sdk.Msg{collection.NewMsgIssueFT(masterAddress, masterAddress, contractID, "ft", "{}", sdk.NewInt(1),
			sdk.NewInt(8), true)},
		service.BroadcastBlock,
	)
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", types.FailedTxError{Tx: res}
	}
	for _, attribute := range res.Logs[0].Events[0].Attributes {
		if attribute.Key == "token_id" {
			return attribute.Value, nil
		}
	}
	return "", types.NoContractIDError{Tx: res}
}

func IssueNFT(linkService *service.LinkService, txBuilder transaction.TxBuilderWithoutKeybase,
	masterKeyWallet *wallet.KeyWallet, contractID string) (string, error) {
	masterAddress := masterKeyWallet.Address()
	account, err := linkService.GetAccount(masterAddress.String())
	if err != nil {
		return "", err
	}
	res, err := linkService.BroadcastMsgs(
		txBuilder,
		account,
		masterKeyWallet,
		[]sdk.Msg{collection.NewMsgIssueNFT(masterAddress, contractID, "nft", "{}")},
		service.BroadcastBlock,
	)
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", types.FailedTxError{Tx: res}
	}
	for _, attribute := range res.Logs[0].Events[0].Attributes {
		if attribute.Key == "token_type" {
			return attribute.Value, nil
		}
	}
	return "", types.NoContractIDError{Tx: res}
}
