package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/codec"
	codectypes "github.com/Finschia/finschia-sdk/codec/types"
	types2 "github.com/Finschia/finschia-sdk/crypto/types"
	"github.com/Finschia/finschia-sdk/simapp"
	"github.com/Finschia/finschia-sdk/store"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/tx/signing"
	authkeeper "github.com/Finschia/finschia-sdk/x/auth/keeper"
	xauthsigning "github.com/Finschia/finschia-sdk/x/auth/signing"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	banktypes "github.com/Finschia/finschia-sdk/x/bank/types"
	bankpluskeeper "github.com/Finschia/finschia-sdk/x/bankplus/keeper"
	feegrant "github.com/Finschia/finschia-sdk/x/feegrant"
	feegrantkeeper "github.com/Finschia/finschia-sdk/x/feegrant/keeper"
	minttypes "github.com/Finschia/finschia-sdk/x/mint/types"
	"github.com/Finschia/finschia-sdk/x/zkauth/keeper"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
)

type TestApp struct {
	Simapp          *simapp.SimApp
	ZKAuthKeeper    *keeper.Keeper
	AccountKeeper   authkeeper.AccountKeeper
	BankKeeper      bankpluskeeper.Keeper
	FeeGrantKeeper  feegrantkeeper.Keeper
	SignModeHandler xauthsigning.SignModeHandler
	Ctx             sdk.Context
	ClientCtx       client.Context
	TxBuilder       client.TxBuilder
}

func ZkAuthKeeper(t testing.TB) TestApp {
	const checkTx = false
	app := simapp.Setup(checkTx)
	maccPerms := simapp.GetMaccPerms()
	appCodec := simapp.MakeTestEncodingConfig().Marshaler

	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	var verificationKey = []byte("{\n \"protocol\": \"groth16\",\n \"curve\": \"bn128\",\n \"nPublic\": 1,\n \"vk_alpha_1\": [\n  \"20491192805390485299153009773594534940189261866228447918068658471970481763042\",\n  \"9383485363053290200918347156157836566562967994039712273449902621266178545958\",\n  \"1\"\n ],\n \"vk_beta_2\": [\n  [\n   \"6375614351688725206403948262868962793625744043794305715222011528459656738731\",\n   \"4252822878758300859123897981450591353533073413197771768651442665752259397132\"\n  ],\n  [\n   \"10505242626370262277552901082094356697409835680220590971873171140371331206856\",\n   \"21847035105528745403288232691147584728191162732299865338377159692350059136679\"\n  ],\n  [\n   \"1\",\n   \"0\"\n  ]\n ],\n \"vk_gamma_2\": [\n  [\n   \"10857046999023057135944570762232829481370756359578518086990519993285655852781\",\n   \"11559732032986387107991004021392285783925812861821192530917403151452391805634\"\n  ],\n  [\n   \"8495653923123431417604973247489272438418190587263600148770280649306958101930\",\n   \"4082367875863433681332203403145435568316851327593401208105741076214120093531\"\n  ],\n  [\n   \"1\",\n   \"0\"\n  ]\n ],\n \"vk_delta_2\": [\n  [\n   \"21349319915249622662700217004338779716430783387183352766870647565870141979289\",\n   \"8213816744021090866451311756048660670381089332123677295675725952502733471420\"\n  ],\n  [\n   \"4787213629490370557685854255230879988945206163033639129474026644007741911075\",\n   \"20003855859301921415178037270191878217707285640767940877063768682564788786247\"\n  ],\n  [\n   \"1\",\n   \"0\"\n  ]\n ],\n \"vk_alphabeta_12\": [\n  [\n   [\n    \"2029413683389138792403550203267699914886160938906632433982220835551125967885\",\n    \"21072700047562757817161031222997517981543347628379360635925549008442030252106\"\n   ],\n   [\n    \"5940354580057074848093997050200682056184807770593307860589430076672439820312\",\n    \"12156638873931618554171829126792193045421052652279363021382169897324752428276\"\n   ],\n   [\n    \"7898200236362823042373859371574133993780991612861777490112507062703164551277\",\n    \"7074218545237549455313236346927434013100842096812539264420499035217050630853\"\n   ]\n  ],\n  [\n   [\n    \"7077479683546002997211712695946002074877511277312570035766170199895071832130\",\n    \"10093483419865920389913245021038182291233451549023025229112148274109565435465\"\n   ],\n   [\n    \"4595479056700221319381530156280926371456704509942304414423590385166031118820\",\n    \"19831328484489333784475432780421641293929726139240675179672856274388269393268\"\n   ],\n   [\n    \"11934129596455521040620786944827826205713621633706285934057045369193958244500\",\n    \"8037395052364110730298837004334506829870972346962140206007064471173334027475\"\n   ]\n  ]\n ],\n \"IC\": [\n  [\n   \"801233197807402683764630185033839955156034586542543249813920835808534245147\",\n   \"13286420793149616228297035344471157585445615731792629462934831296345279687002\",\n   \"1\"\n  ],\n  [\n   \"17608180544527043978731301492557909061209088433544687588079992534282036547698\",\n   \"11240405619785894451348234456278767489162139374206168239508590931049712428392\",\n   \"1\"\n  ]\n ]\n}")

	jwKsMap := types.NewJWKs()
	jwKsMap.AddJWK(&types.JWK{
		Kty: "RSA",
		E:   "AQAB",
		N:   "q0CrF3x3aYsjr0YOLMOAhEGMvyFp6o4RqyEdUrnTDYkhZbcud-fJEQafCTnjS9QHN1IjpuK6gpx5i3-Z63vRjs5EQX7lP1jG8Qg-CnBdTTLw4uJi7RmmlKPsYaO1DbNkFO2uEN62sOOzmJCh1od3CZXI1UYH5cvZ_sLJaN2A4TwvUTU3aXlXbUNJz_Hy3l0q1Jjta75NrJtJ7Pfj9tVXs8qXp15tZXrnbaM-AI0puswt35VsQbmLwUovFFGeToo5q2c_c1xYnV5uQYMadANekGPRFPM9JZpSSIvH0Lv_f15V2zRqmIgX7a3RcmTnr3-w3QNQTogdy-MogxPUdRbxow",
		Alg: "RS256",
		Kid: "55c188a83546fc188e51576ba72836e0600e8b73",
	})
	jwKsMap.AddJWK(&types.JWK{
		N:   "pOpd5-7RpMvcfBcSjqlTNYjGg3YRwYRV9T9k7eDOEWgMBQEs6ii3cjcuoa1oD6N48QJmcNvAme_ud985DV2mQpOaCUy22MVRKI8DHxAKGWzZO5yzn6otsN9Vy0vOEO_I-vnmrO1-1ONFuH2zieziaXCUVh9087dRkM9qaQYt6QJhMmiNpyrbods6AsU8N1jeAQl31ovHWGGk8axXNmwbx3dDZQhx-t9ZD31oF-usPhFZtM92mxgehDqi2kpvFmM0nzSVgPrOXlbDb9ztg8lclxKwnT1EtcwHUq4FeuOPQMtZ2WehrY10OvsqS5ml3mxXUQEXrtYfa5V1v4o3rWx9Ow",
		Kty: "RSA",
		Alg: "RS256",
		E:   "AQAB",
		Kid: "6f9777a685907798ef794062c00b65d66c240b1b",
	})

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		jwKsMap,
		types.NewZKAuthVerifier(verificationKey),
		app.MsgServiceRouter(),
	)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	encodingConfig := simapp.MakeTestEncodingConfig()
	encodingConfig.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)
	testdata.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	clientCtx := client.Context{}.WithTxConfig(encodingConfig.TxConfig)
	authKeeper := authkeeper.NewAccountKeeper(
		appCodec, app.GetKey(authtypes.StoreKey), app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms,
	)
	feeGrantKeeper := feegrantkeeper.NewKeeper(appCodec, app.GetKey(feegrant.StoreKey), authKeeper)
	bankKeeper := bankpluskeeper.NewBaseKeeper(
		appCodec, app.GetKey(banktypes.StoreKey), authKeeper, app.GetSubspace(banktypes.ModuleName), app.BlockedAddrs(), false)

	testApp := TestApp{
		Simapp:          app,
		ZKAuthKeeper:    k,
		AccountKeeper:   authKeeper,
		BankKeeper:      bankKeeper,
		FeeGrantKeeper:  feeGrantKeeper,
		SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
		Ctx:             ctx,
		ClientCtx:       clientCtx,
		TxBuilder:       clientCtx.TxConfig.NewTxBuilder(),
	}

	return testApp
}

func (t *TestApp) CreateTestAccounts(numAcc int) ([]authtypes.AccountI, error) {
	var accounts []authtypes.AccountI

	for i := 0; i < numAcc; i++ {
		_, _, addr := testdata.KeyTestPubAddr()
		acc, err := t.addAccount(addr, i)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (t *TestApp) addAccount(accAddr sdk.AccAddress, accNum int) (authtypes.AccountI, error) {
	acc := t.Simapp.AccountKeeper.NewAccountWithAddress(t.Ctx, accAddr)
	if err := acc.SetAccountNumber(uint64(accNum)); err != nil {
		return nil, err
	}

	t.Simapp.AccountKeeper.SetAccount(t.Ctx, acc)
	someCoins := sdk.Coins{sdk.NewInt64Coin("cony", 10000000)}
	if err := t.Simapp.BankKeeper.MintCoins(t.Ctx, minttypes.ModuleName, someCoins); err != nil {
		return nil, err
	}

	if err := t.Simapp.BankKeeper.SendCoinsFromModuleToAccount(t.Ctx, minttypes.ModuleName, accAddr, someCoins); err != nil {
		return nil, err
	}

	return acc, nil
}

func (t *TestApp) AddTestAccounts(addrs []string) ([]authtypes.AccountI, error) {
	accounts := make([]authtypes.AccountI, 0)

	for i, addrStr := range addrs {
		addr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			return nil, err
		}

		acc, err := t.addAccount(addr, i)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)
	}

	return accounts, nil
}

func (t *TestApp) CreateTestTx(pubs []types2.PubKey, accSeqs []uint64) (xauthsigning.Tx, error) {
	sigsV2 := make([]signing.SignatureV2, 0)
	for i, pub := range pubs {
		sigV2 := signing.SignatureV2{
			PubKey: pub,
			Data: &signing.SingleSignatureData{
				SignMode:  t.ClientCtx.TxConfig.SignModeHandler().DefaultMode(),
				Signature: nil,
			},
			Sequence: accSeqs[i],
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err := t.TxBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	return t.TxBuilder.GetTx(), nil
}
