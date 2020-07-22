package scenario

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/wallet"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type TxTokenScenario struct {
	Info
}

func (s *TxTokenScenario) GenerateStateSettingMsgs(masterKeyWallet *wallet.KeyWallet,
	hdWallet *wallet.HDWallet, scenarioParams []string) ([]sdk.Msg, map[string]string, error) {
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
	return append(registerMsgs, grantPermMsgs...), map[string]string{"token_contract_id": contractID}, nil
}

func (s *TxTokenScenario) GenerateTarget(keyWallet *wallet.KeyWallet, walletIndex int) (*[]*vegeta.Target, int, error) {
	target, err := BuildTxTarget(s.Info, keyWallet, walletIndex, []string{"MsgMint", "MsgTransfer",
		"MsgGrantPermission", "MsgModifyToken", "MsgBurn"}, s.scenarioParams)
	targets := []*vegeta.Target{target}
	return &targets, len(targets), err
}
