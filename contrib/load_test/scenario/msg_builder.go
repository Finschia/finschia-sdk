package scenario

import (
	"fmt"
	"math"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/types"
	linktypes "github.com/line/link/types"
	acc "github.com/line/link/x/account"
	"github.com/line/link/x/coin"
	"github.com/line/link/x/collection"
	"github.com/line/link/x/token"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type MsgBuilder struct {
	fromAddr       sdk.AccAddress
	walletIndex    int
	coinName       string
	scenarioParams []string

	tokenContractID      string
	collectionContractID string
	ftTokenID            string
	nftTokenType         string
	numNFTPerUser        int
	nftTokenIDs          []string
	nftOffset            int
	modifyCounter        int
	childNFT             int

	handlers map[string]func() (sdk.Msg, error)
}

func NewMsgBuilder(fromAddr sdk.AccAddress, info Info, walletIndex int, scenarioParams []string) (*MsgBuilder, error) {
	var numNFTPerUser int
	var err error

	if val, ok := info.stateParams["num_nft_per_user"]; ok {
		numNFTPerUser, err = strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
	}

	if numNFTPerUser*(walletIndex+1) > math.MaxInt32 {
		return nil, types.NFTTokenIDOverFlowError{}
	}
	nftTokenIDs := make([]string, numNFTPerUser)
	for i := 0; i < numNFTPerUser; i++ {
		nftTokenIDs[i] = fmt.Sprintf("%s%08x", info.stateParams["nft_token_type"], numNFTPerUser*walletIndex+i+1)
	}

	m := &MsgBuilder{
		fromAddr:             fromAddr,
		walletIndex:          walletIndex,
		coinName:             info.config.CoinName,
		scenarioParams:       scenarioParams,
		tokenContractID:      info.stateParams["token_contract_id"],
		collectionContractID: info.stateParams["collection_contract_id"],
		ftTokenID:            info.stateParams["ft_token_id"],
		nftTokenType:         info.stateParams["nft_token_type"],
		numNFTPerUser:        numNFTPerUser,
		nftTokenIDs:          nftTokenIDs,
		nftOffset:            0,
		modifyCounter:        0,
		handlers:             make(map[string]func() (sdk.Msg, error)),
	}

	m.handlers["MsgEmpty"] = m.newMsgEmpty
	m.handlers["MsgSend"] = m.newMsgSend
	m.handlers["MsgIssue"] = m.newMsgIssue
	m.handlers["MsgMint"] = m.newMsgMint
	m.handlers["MsgTransfer"] = m.newMsgTransfer
	m.handlers["MsgModifyToken"] = m.newMsgModifyToken
	m.handlers["MsgModifyTokenName"] = m.newMsgModifyTokenName
	m.handlers["MsgModifyTokenURI"] = m.newMsgModifyTokenURI
	m.handlers["MsgBurn"] = m.newMsgBurn
	m.handlers["MsgCreateCollection"] = m.newMsgCreateCollection
	m.handlers["MsgApprove"] = m.newMsgApprove
	m.handlers["MsgIssueFT"] = m.newMsgIssueFT
	m.handlers["MsgMintFT"] = m.newMsgMintFT
	m.handlers["MsgTransferFT"] = m.newMsgTransferFT
	m.handlers["MsgBurnFT"] = m.newMsgBurnFT
	m.handlers["MsgModifyCollection"] = m.newMsgModifyCollection
	m.handlers["MsgIssueNFT"] = m.newMsgIssueNFT
	m.handlers["MsgMintNFT"] = m.newMsgMintNFT
	m.handlers["MsgMintOneNFT"] = m.newMsgMintOneNFT
	m.handlers["MsgMintFiveNFT"] = m.newMsgMintFiveNFTs
	m.handlers["MsgAttach"] = m.newMsgAttach
	m.handlers["MsgDetach"] = m.newMsgDetach
	m.handlers["MsgTransferNFT"] = m.newMsgTransferNFT
	m.handlers["MsgMultiTransferNFT"] = m.newMsgMultiTransferNFT
	m.handlers["MsgBurnNFT"] = m.newMsgBurnNFT
	m.handlers["MsgGrantPermission"] = m.newMsgGrantPermission

	return m, nil
}

func (m *MsgBuilder) GetHandler(msgType string) (func() (sdk.Msg, error), error) {
	if _, ok := m.handlers[msgType]; !ok {
		return nil, types.NoMsgBuildHandler{MsgType: msgType}
	}
	return m.handlers[msgType], nil
}

func (m *MsgBuilder) newMsgEmpty() (sdk.Msg, error) {
	return acc.NewMsgEmpty(m.fromAddr), nil
}

func (m *MsgBuilder) newMsgSend() (sdk.Msg, error) {
	return coin.NewMsgSend(
		m.fromAddr,
		secp256k1.GenPrivKey().PubKey().Address().Bytes(),
		sdk.NewCoins(sdk.NewCoin(m.coinName, sdk.NewInt(1))),
	), nil
}

func (m *MsgBuilder) newMsgIssue() (sdk.Msg, error) {
	return token.NewMsgIssue(
		m.fromAddr,
		secp256k1.GenPrivKey().PubKey().Address().Bytes(),
		"token",
		"TOK",
		"{}",
		"uri",
		sdk.NewInt(1),
		sdk.NewInt(8),
		true,
	), nil
}

func (m *MsgBuilder) newMsgMint() (sdk.Msg, error) {
	return token.NewMsgMint(
		m.fromAddr,
		m.tokenContractID,
		m.fromAddr,
		sdk.NewInt(2),
	), nil
}

func (m *MsgBuilder) newMsgTransfer() (sdk.Msg, error) {
	return token.NewMsgTransfer(
		m.fromAddr,
		secp256k1.GenPrivKey().PubKey().Address().Bytes(),
		m.tokenContractID,
		sdk.NewInt(1),
	), nil
}

func (m *MsgBuilder) newMsgModifyToken() (sdk.Msg, error) {
	return m.newMsgModifyTokenName()
}
func (m *MsgBuilder) newMsgModifyTokenName() (sdk.Msg, error) {
	m.modifyCounter++
	return token.NewMsgModify(
		m.fromAddr,
		m.tokenContractID,
		linktypes.NewChangesWithMap(map[string]string{"name": fmt.Sprintf("token%d-%d", m.walletIndex,
			m.modifyCounter)}),
	), nil
}

func (m *MsgBuilder) newMsgModifyTokenURI() (sdk.Msg, error) {
	m.modifyCounter++
	return token.NewMsgModify(
		m.fromAddr,
		m.tokenContractID,
		linktypes.NewChangesWithMap(map[string]string{"img_uri": fmt.Sprintf("uri%d-%d", m.walletIndex,
			m.modifyCounter)}),
	), nil
}
func (m *MsgBuilder) newMsgBurn() (sdk.Msg, error) {
	return token.NewMsgBurn(
		m.fromAddr,
		m.tokenContractID,
		sdk.NewInt(1),
	), nil
}

func (m *MsgBuilder) newMsgCreateCollection() (sdk.Msg, error) {
	return collection.NewMsgCreateCollection(
		m.fromAddr,
		"name",
		"{}",
		"uri",
	), nil
}

func (m *MsgBuilder) newMsgApprove() (sdk.Msg, error) {
	return collection.NewMsgApprove(
		m.fromAddr,
		m.collectionContractID,
		secp256k1.GenPrivKey().PubKey().Address().Bytes(),
	), nil
}

func (m *MsgBuilder) newMsgIssueFT() (sdk.Msg, error) {
	return collection.NewMsgIssueFT(
		m.fromAddr,
		secp256k1.GenPrivKey().PubKey().Address().Bytes(),
		m.collectionContractID,
		"collection",
		"{}",
		sdk.NewInt(1),
		sdk.NewInt(8),
		true,
	), nil
}

func (m *MsgBuilder) newMsgMintFT() (sdk.Msg, error) {
	return collection.NewMsgMintFT(
		m.fromAddr,
		m.collectionContractID,
		m.fromAddr,
		collection.NewCoin(m.ftTokenID, sdk.NewInt(4)),
	), nil
}

func (m *MsgBuilder) newMsgTransferFT() (sdk.Msg, error) {
	return collection.NewMsgTransferFT(
		m.fromAddr,
		m.collectionContractID,
		secp256k1.GenPrivKey().PubKey().Address().Bytes(),
		collection.NewCoin(m.ftTokenID, sdk.NewInt(1)),
	), nil
}

func (m *MsgBuilder) newMsgBurnFT() (sdk.Msg, error) {
	return collection.NewMsgBurnFT(
		m.fromAddr,
		m.collectionContractID,
		collection.NewCoin(m.ftTokenID, sdk.NewInt(1)),
	), nil
}

func (m *MsgBuilder) newMsgModifyCollection() (sdk.Msg, error) {
	m.modifyCounter++
	return collection.NewMsgModify(
		m.fromAddr,
		m.collectionContractID,
		m.ftTokenID[:8],
		m.ftTokenID[8:],
		linktypes.NewChangesWithMap(map[string]string{"name": fmt.Sprintf("name%d-%d", m.walletIndex,
			m.modifyCounter)}),
	), nil
}

func (m *MsgBuilder) newMsgIssueNFT() (sdk.Msg, error) {
	return collection.NewMsgIssueNFT(
		m.fromAddr,
		m.collectionContractID,
		"ft",
		"{}",
	), nil
}

func (m *MsgBuilder) newMsgMintNFT() (sdk.Msg, error) {
	numNFT, err := strconv.Atoi(m.scenarioParams[0])
	if err != nil {
		return nil, err
	}
	return m.newMsgMintNFTs(numNFT)
}

func (m *MsgBuilder) newMsgMintOneNFT() (sdk.Msg, error) {
	return m.newMsgMintNFTs(1)
}

func (m *MsgBuilder) newMsgMintFiveNFTs() (sdk.Msg, error) {
	return m.newMsgMintNFTs(5)
}

func (m *MsgBuilder) newMsgMintNFTs(numNFT int) (sdk.Msg, error) {
	params := make([]collection.MintNFTParam, 0)
	for j := 0; j < numNFT; j++ {
		params = append(params, collection.NewMintNFTParam("name", "{}", m.nftTokenType))
	}
	return collection.NewMsgMintNFT(
		m.fromAddr,
		m.collectionContractID,
		secp256k1.GenPrivKey().PubKey().Address().Bytes(),
		params...,
	), nil
}

func (m *MsgBuilder) newMsgAttach() (sdk.Msg, error) {
	if m.nftOffset+1 >= m.numNFTPerUser {
		return nil, types.OutOfNFTError{NumNFTPerUser: m.numNFTPerUser, NFTOffset: m.nftOffset}
	}
	m.childNFT = m.nftOffset + 1
	return collection.NewMsgAttach(
		m.fromAddr,
		m.collectionContractID,
		m.nftTokenIDs[m.nftOffset],
		m.nftTokenIDs[m.childNFT],
	), nil
}

func (m *MsgBuilder) newMsgDetach() (sdk.Msg, error) {
	return collection.NewMsgDetach(
		m.fromAddr,
		m.collectionContractID,
		m.nftTokenIDs[m.childNFT],
	), nil
}

func (m *MsgBuilder) newMsgTransferNFT() (sdk.Msg, error) {
	if m.nftOffset >= m.numNFTPerUser {
		return nil, types.OutOfNFTError{NumNFTPerUser: m.numNFTPerUser, NFTOffset: m.nftOffset}
	}
	msg := collection.NewMsgTransferNFT(
		m.fromAddr,
		m.collectionContractID,
		secp256k1.GenPrivKey().PubKey().Address().Bytes(),
		m.nftTokenIDs[m.nftOffset],
	)
	m.nftOffset++
	return msg, nil
}

func (m *MsgBuilder) newMsgMultiTransferNFT() (sdk.Msg, error) {
	if m.nftOffset+5 > m.numNFTPerUser {
		return nil, types.OutOfNFTError{NumNFTPerUser: m.numNFTPerUser, NFTOffset: m.nftOffset}
	}
	msg := collection.NewMsgTransferNFT(
		m.fromAddr,
		m.collectionContractID,
		secp256k1.GenPrivKey().PubKey().Address().Bytes(),
		m.nftTokenIDs[m.nftOffset:m.nftOffset+5]...,
	)
	m.nftOffset += 5
	return msg, nil
}

func (m *MsgBuilder) newMsgBurnNFT() (sdk.Msg, error) {
	if m.nftOffset >= m.numNFTPerUser {
		return nil, types.OutOfNFTError{NumNFTPerUser: m.numNFTPerUser, NFTOffset: m.nftOffset}
	}
	msg := collection.NewMsgBurnNFT(
		m.fromAddr,
		m.collectionContractID,
		m.nftTokenIDs[m.nftOffset],
	)
	m.nftOffset++
	return msg, nil
}

func (m *MsgBuilder) newMsgGrantPermission() (sdk.Msg, error) {
	return token.NewMsgGrantPermission(
		m.fromAddr,
		m.tokenContractID,
		secp256k1.GenPrivKey().PubKey().Address().Bytes(),
		"modify",
	), nil
}
