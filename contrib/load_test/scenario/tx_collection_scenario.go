package scenario

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/transaction"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	linktypes "github.com/line/link/types"
	"github.com/line/link/x/collection"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type TxCollectionScenario struct {
	linkService   *service.LinkService
	targetBuilder *TargetBuilder
	txBuilder     transaction.TxBuilderWithoutKeybase
	config        types.Config
	params        map[string]string
}

func (s *TxCollectionScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
	hdWallet *wallet.HDWallet) ([]sdk.Msg, map[string]string, error) {
	masterAddress := masterKeyWallet.Address()
	contractID, err := CreateCollection(s.linkService, s.txBuilder, masterKeyWallet)
	if err != nil {
		return nil, nil, err
	}
	ftTokenID, err := IssueFT(s.linkService, s.txBuilder, masterKeyWallet, contractID)
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
		[]string{"mint", "burn", "modify"})
	if err != nil {
		return nil, nil, err
	}
	mintNFTMsgs, err := GenerateMintNFTMsgs(masterAddress, hdWallet, s.config, contractID, nftTokenType, 2)
	if err != nil {
		return nil, nil, err
	}
	msgs = append(msgs, grantPermMsgs...)
	msgs = append(msgs, mintNFTMsgs...)

	return msgs, map[string]string{"contract_id": contractID, "ft_token_id": ftTokenID,
		"nft_token_type": nftTokenType}, nil
}

func (s *TxCollectionScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	account, err := s.linkService.GetAccount(keyWallet.Address().String())
	if err != nil {
		return nil, 0, err
	}

	numMsgIncrement := 8
	repeatCount := (s.config.MsgsPerTxLoadTest + numMsgIncrement - 1) / numMsgIncrement // round up
	from := account.Address
	msgs := make([]sdk.Msg, numMsgIncrement*repeatCount)
	contractID := s.params["contract_id"]
	ftTokenID := s.params["ft_token_id"]
	nftTokenType := s.params["nft_token_type"]
	nftTokenID1 := fmt.Sprintf("%s%08x", nftTokenType, 2*walletIndex+1)
	nftTokenID2 := fmt.Sprintf("%s%08x", nftTokenType, 2*walletIndex+2)
	for i := 0; i < repeatCount; i++ {
		to := secp256k1.GenPrivKey().PubKey().Address().Bytes()
		msgs[numMsgIncrement*i] = collection.NewMsgMintFT(from, contractID, from, collection.NewCoin(ftTokenID,
			sdk.NewInt(2)))
		msgs[numMsgIncrement*i+1] = collection.NewMsgTransferFT(from, contractID, to, collection.NewCoin(ftTokenID,
			sdk.NewInt(1)))
		msgs[numMsgIncrement*i+2] = collection.NewMsgModify(from, contractID, ftTokenID[:8], ftTokenID[8:],
			linktypes.NewChangesWithMap(map[string]string{"name": "new_name"}))
		msgs[numMsgIncrement*i+3] = collection.NewMsgBurnFT(from, contractID, collection.NewCoin(ftTokenID,
			sdk.NewInt(1)))
		msgs[numMsgIncrement*i+4] = collection.NewMsgAttach(from, contractID, nftTokenID1, nftTokenID2)
		msgs[numMsgIncrement*i+5] = collection.NewMsgDetach(from, contractID, nftTokenID2)
		msgs[numMsgIncrement*i+6] = collection.NewMsgTransferNFT(from, contractID, to, nftTokenID2)
		msgs[numMsgIncrement*i+7] = collection.NewMsgBurnNFT(from, contractID, nftTokenID1)
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
