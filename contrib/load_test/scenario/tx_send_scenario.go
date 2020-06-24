package scenario

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/transaction"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	"github.com/line/link/x/coin"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type TxSendScenario struct {
	linkService   *service.LinkService
	targetBuilder *TargetBuilder
	txBuilder     transaction.TxBuilderWithoutKeybase
	config        types.Config
}

func (s *TxSendScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
	hdWallet *wallet.HDWallet) ([]sdk.Msg, map[string]string, error) {
	msgs, err := GenerateRegisterAccountMsgs(masterKeyWallet.Address(), hdWallet, s.config)
	return msgs, nil, err
}

func (s *TxSendScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	account, err := s.linkService.GetAccount(keyWallet.Address().String())
	if err != nil {
		return nil, 0, err
	}

	from := account.Address
	coins := sdk.NewCoins(sdk.NewCoin(s.config.CoinName, sdk.NewInt(1)))
	msgs := make([]sdk.Msg, s.config.MsgsPerTxLoadTest)
	for i := 0; i < s.config.MsgsPerTxLoadTest; i++ {
		to := secp256k1.GenPrivKey().PubKey().Address().Bytes()
		msgs[i] = coin.NewMsgSend(from, to, coins)
	}

	stdTx, err := s.txBuilder.WithAccountNumber(account.AccountNumber).WithSequence(account.Sequence).
		BuildAndSign(keyWallet.PrivateKey(), msgs)
	if err != nil {
		return nil, 0, err
	}

	target, err := s.targetBuilder.MakeTxTarget(stdTx, service.BroadcastSync)
	targets := []*vegeta.Target{target}
	return &targets, len(targets), err
}
