package keeper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Finschia/finschia-sdk/baseapp"
	"github.com/Finschia/finschia-sdk/codec"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
	"github.com/Finschia/ostracon/libs/log"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	jwks     []types.JWK
	lock     sync.RWMutex
	router   *baseapp.MsgServiceRouter
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	router *baseapp.MsgServiceRouter,
) *Keeper {

	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		jwks:     nil,
		router:   router,
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) DispatchMsgs(ctx sdk.Context, msgs []sdk.Msg) ([][]byte, error) {
	results := make([][]byte, len(msgs))

	for i, msg := range msgs {
		signers := msg.GetSigners()
		if len(signers) != 1 {
			return nil, sdkerrors.ErrInvalidRequest.Wrap("authorization can be given to msg with only one signer")
		}
		if err := msg.ValidateBasic(); err != nil {
			return nil, err
		}

		handler := k.router.Handler(msg)
		if handler == nil {
			return nil, sdkerrors.ErrUnknownRequest.Wrapf("unrecognized message route: %s", sdk.MsgTypeURL(msg))
		}

		msgResp, err := handler(ctx, msg)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "failed to execute message; message %v", msg)
		}

		results[i] = msgResp.Data

		events := msgResp.Events
		sdkEvents := make([]sdk.Event, 0, len(events))
		for _, event := range events {
			e := event
			e.Attributes = append(e.Attributes, abci.EventAttribute{Key: []byte("zkauth_msg_index"), Value: []byte(strconv.Itoa(i))})

			sdkEvents = append(sdkEvents, sdk.Event(e))
		}

		ctx.EventManager().EmitEvents(sdkEvents)
	}

	return results, nil
}

func (k *Keeper) LoopJWK(ctx sdk.Context) {
	logger := k.Logger(ctx)
	var defaultZKAuthOAuthProviders = [1]types.OidcProvider{types.Google}
	var fetchIntervals uint64
	go func() {
		for {
			select {
			case <-ctx.Context().Done():
				return
			default:
				for _, name := range defaultZKAuthOAuthProviders {
					provider := types.GetConfig(name)
					fetchIntervals = provider.FetchIntervals

					err := k.FetchJWK(provider.JwkEndpoint, name)
					if err != nil {
						time.Sleep(time.Duration(fetchIntervals) * time.Second)
						logger.Error(fmt.Sprintf("%s", err))
						continue
					}
				}
				time.Sleep(time.Duration(fetchIntervals) * time.Second)
			}
		}
	}()
}

func (k *Keeper) ParseJWKs(byteArray []byte) (jwks []types.JWK, err error) {
	var data map[string]interface{}
	err = json.Unmarshal(byteArray, &data)
	if err != nil {
		return jwks, err
	}

	for _, v := range data["keys"].([]interface{}) {
		var jwk types.JWK

		b, err := json.Marshal(v)
		if err != nil {
			return jwks, err
		}

		if err := json.Unmarshal(b, &jwk); err != nil {
			return jwks, err
		}

		jwks = append(jwks, jwk)
	}

	return jwks, nil
}

func (k *Keeper) FetchJWK(endpoint string, name types.OidcProvider) error {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	dataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	jwks, err := k.ParseJWKs(dataBytes)
	if err != nil {
		return err
	}

	k.SetJWKs(jwks)
	resp.Body.Close()

	return nil
}

func (k *Keeper) GetJWKs() []types.JWK {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.jwks
}

func (k *Keeper) SetJWKs(jwks []types.JWK) {
	k.lock.Lock()
	defer k.lock.Unlock()

	k.jwks = jwks
}

func (k *Keeper) GetJWK(kid string) *types.JWK {
	for _, v := range k.GetJWKs() {
		if v.Kid == kid {
			return &v
		}
	}
	return nil
}

func (k *Keeper) CreateJWKFileName(name types.OidcProvider) string {
	return fmt.Sprintf("jwk-%s.json", name)
}
