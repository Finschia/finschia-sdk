package keeper_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/simapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	banktypes "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/zkauth/testutil"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
)

const testData = `{
	"keys": [
	  {
		"kid": "bdc4e109815f469460e63d34cd684215148d7b59",
		"e": "AQAB",
		"kty": "RSA",
		"alg": "RS256",
		"n": "v3dZL2R2PuebbAChYXKVW6R-FJDUVmZ8TyVMWH0-VpVjFYZvy7BZaE5ApLWc3UhpXug6r6230AJI0ow5yePnqmZnI5qckxz0br0Fj27Zdg-X4PWN95gdk6fpI4JwNmZFsgiWzmDiP118j8jIxMNBiIVPT7RyykhAZeNnGC2kDU-81iop850K205EwfSi_TBT6HCbRj_TSQ2oJfIXDPX8s7Kg4PRjDOHt3D8CiqsIWbxSkRRuTiU_1Ahsbuc3d9hkD1rOOThVT6T7LVZT710WtPa1QbKUgGIu2pmiPo0BCdnbqozsRVOwY901R77VlVwpTuGonPZuyO1B2FgGuYgotw",
		"use": "sig"
	  },
	  {
		"kty": "RSA",
		"e": "AQAB",
		"kid": "ed806f1842b588054b18b669dd1a09a4f367afc4",
		"n": "rH3Q5NY6MAeaE8NuSw7Rw2Cc1e_j-kUS044tu-WcmTFzBKTuKvIlgj5w0SlSbiVl81zBtetQFtuwkMzWgnCks-2-Fwpoy__2NUouUgLtIggAVEyOGgPLfyaswtkSmZsUmWWg9J8CgMUdoXFkbZAPladDcmSqiXJ7cp9nvro6f4sjfrGDYz5_-SNz1AQEGbvcTh9EeZkvKPrmnV3YER95bJsgkHmNJVkQ6LcWtLyKhSGQGRMeTYaXDajc2KrKT3net7qNhbAm7KpWddbtR5l6A0TRCrAMoV2M68_GLRF24acj3UO5RW0SkuaBTZS4KQpyoyABCAtjLSr-3RY6WR9npw",
		"alg": "RS256",
		"use": "sig"
	  },
	  {
		"n": "q0CrF3x3aYsjr0YOLMOAhEGMvyFp6o4RqyEdUrnTDYkhZbcud-fJEQafCTnjS9QHN1IjpuK6gpx5i3-Z63vRjs5EQX7lP1jG8Qg-CnBdTTLw4uJi7RmmlKPsYaO1DbNkFO2uEN62sOOzmJCh1od3CZXI1UYH5cvZ_sLJaN2A4TwvUTU3aXlXbUNJz_Hy3l0q1Jjta75NrJtJ7Pfj9tVXs8qXp15tZXrnbaM-AI0puswt35VsQbmLwUovFFGeToo5q2c_c1xYnV5uQYMadANekGPRFPM9JZpSSIvH0Lv_f15V2zRqmIgX7a3RcmTnr3-w3QNQTogdy-MogxPUdRbxow",
		"e": "AQAB",
		"kid": "55c188a83546fc188e51576ba72836e0600e8b73",
		"kty": "RSA",
		"use": "sig",
		"alg": "RS256"
	  }
	]
  }`

func mockHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonData := testData
	w.Write([]byte(jsonData))
}

func TestFetchJwk(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer server.Close()
	testApp := testutil.ZkAuthKeeper(t)
	k := testApp.ZKAuthKeeper
	ctx := testApp.Ctx
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tempDir, err := os.MkdirTemp("", types.StoreKey)
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	k.FetchJWK(ctx.WithContext(timeoutCtx))
	<-timeoutCtx.Done()

	require.True(t, k.GetJWKSize() == 3)

	var expectedObj types.JWKs
	err = json.Unmarshal([]byte(testData), &expectedObj)
	require.NoError(t, err)
}

func TestDispatchMsgs(t *testing.T) {
	testApp := testutil.ZkAuthKeeper(t)
	app, k, ctx := testApp.Simapp, testApp.ZKAuthKeeper, testApp.Ctx

	addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(100))
	fromAddr := addrs[0]
	toAddr := addrs[1]

	newCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 5))

	bankMsg := banktypes.MsgSend{
		Amount:      newCoins,
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
	}

	zkAuthSig := types.ZKAuthSignature{}

	msgs := types.NewMsgExecution([]sdk.Msg{&bankMsg}, zkAuthSig)

	execMsgs, err := msgs.GetMessages()
	require.NoError(t, err)
	result, err := k.DispatchMsgs(ctx, execMsgs)

	require.NoError(t, err)
	require.NotNil(t, result)

	fromBalance := app.BankKeeper.GetBalance(ctx, fromAddr, "stake")
	require.True(t, fromBalance.Equal(sdk.NewInt64Coin("stake", 95)))
	toBalance := app.BankKeeper.GetBalance(ctx, toAddr, "stake")
	require.True(t, toBalance.Equal(sdk.NewInt64Coin("stake", 105)))
}
