package scenario

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/transaction"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	linktypes "github.com/line/link/types"
	"github.com/line/link/x/token"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type TxTokenScenario struct {
	linkService   *service.LinkService
	targetBuilder *TargetBuilder
	txBuilder     transaction.TxBuilderWithoutKeybase
	config        types.Config
	params        map[string]string
}

func (s *TxTokenScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
	hdWallet *wallet.HDWallet) ([]sdk.Msg, map[string]string, error) {
	masterAddress := masterKeyWallet.Address()
	contractID, _, err := IssueToken(s.linkService, s.txBuilder, masterKeyWallet)
	if err != nil {
		return nil, nil, err
	}
	registerMsgs, err := GenerateRegisterAccountMsgs(masterAddress, hdWallet, s.config)
	if err != nil {
		return nil, nil, err
	}
	grantPermMsgs, err := GenerateGrantPermissionMsgs(masterAddress, hdWallet, s.config, contractID, "token",
		[]string{"mint", "burn", "modify"})
	if err != nil {
		return nil, nil, err
	}
	return append(registerMsgs, grantPermMsgs...), map[string]string{"contract_id": contractID}, nil
}

func (s *TxTokenScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	account, err := s.linkService.GetAccount(keyWallet.Address().String())
	if err != nil {
		return nil, 0, err
	}

	numMsgIncrement := 6
	repeatCount := (s.config.MsgsPerTxLoadTest + numMsgIncrement - 1) / numMsgIncrement // round up
	from := account.Address

	msgs := make([]sdk.Msg, numMsgIncrement*repeatCount)
	contractID := s.params["contract_id"]
	for i := 0; i < repeatCount; i++ {
		to := secp256k1.GenPrivKey().PubKey().Address().Bytes()
		msgs[numMsgIncrement*i] = token.NewMsgMint(from, contractID, from, sdk.NewInt(2))
		msgs[numMsgIncrement*i+1] = token.NewMsgTransfer(from, to, contractID, sdk.NewInt(1))
		msgs[numMsgIncrement*i+2] = token.NewMsgGrantPermission(from, contractID, to, "mint")
		msgs[numMsgIncrement*i+3] = token.NewMsgRevokePermission(from, contractID, "mint")
		msgs[numMsgIncrement*i+4] = token.NewMsgModify(from, contractID,
			linktypes.NewChangesWithMap(map[string]string{"name": "new_name"}))
		msgs[numMsgIncrement*i+5] = token.NewMsgBurn(from, contractID, sdk.NewInt(1))
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
