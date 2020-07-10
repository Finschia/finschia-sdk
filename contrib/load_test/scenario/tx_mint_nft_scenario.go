package scenario

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/transaction"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	"github.com/line/link/x/collection"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type TxMintNFTScenario struct {
	linkService   *service.LinkService
	targetBuilder *TargetBuilder
	txBuilder     transaction.TxBuilderWithoutKeybase
	config        types.Config
	params        map[string]string
}

func (s *TxMintNFTScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
	hdWallet *wallet.HDWallet) ([]sdk.Msg, map[string]string, error) {
	masterAddress := masterKeyWallet.Address()
	contractID, err := CreateCollection(s.linkService, s.txBuilder, masterKeyWallet)
	if err != nil {
		return nil, nil, err
	}
	nftTokenType, err := IssueNFT(s.linkService, s.txBuilder, masterKeyWallet, contractID)
	if err != nil {
		return nil, nil, err
	}

	msgs, err := GenerateRegisterAccountMsgs(masterAddress, hdWallet, s.config)
	if err != nil {
		return nil, nil, err
	}
	grantPermMsgs, err := GenerateGrantPermissionMsgs(masterAddress, hdWallet, s.config, contractID, "collection",
		[]string{"mint"})
	if err != nil {
		return nil, nil, err
	}
	msgs = append(msgs, grantPermMsgs...)

	return msgs, map[string]string{"contract_id": contractID, "nft_token_type": nftTokenType}, nil
}

func (s *TxMintNFTScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	account, err := s.linkService.GetAccount(keyWallet.Address().String())
	if err != nil {
		return nil, 0, err
	}

	from := account.Address
	msgs := make([]sdk.Msg, s.config.MsgsPerTxLoadTest)
	contractID := s.params["contract_id"]
	nftTokenType := s.params["nft_token_type"]
	for i := 0; i < s.config.MsgsPerTxLoadTest; i++ {
		to := secp256k1.GenPrivKey().PubKey().Address().Bytes()
		msgs[i] = collection.NewMsgMintNFT(from, contractID, to, collection.NewMintNFTParam("name", "{}", nftTokenType))
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
