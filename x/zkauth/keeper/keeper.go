package keeper

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Finschia/ostracon/libs/log"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Finschia/finschia-sdk/baseapp"
	"github.com/Finschia/finschia-sdk/codec"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	jwksMap    *types.JWKsMap       // JWK manager
	zkVerifier types.ZKAuthVerifier // zkp verification key byte
	router     *baseapp.MsgServiceRouter
}

var _ types.ZKAuthKeeper = &Keeper{}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	jwksMap *types.JWKsMap,
	zkVerifier types.ZKAuthVerifier,
	router *baseapp.MsgServiceRouter,
) *Keeper {
	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		jwksMap:    jwksMap,
		zkVerifier: zkVerifier,
		router:     router,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) DispatchMsgs(ctx sdk.Context, msgs []sdk.Msg) ([][]byte, error) {
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

func (k Keeper) FetchJWK(ctx sdk.Context) {
	logger := k.Logger(ctx)
	var defaultZKAuthOAuthProviders = [1]types.OidcProvider{types.Google}
	var fetchIntervals uint64
	go func() {
		for {
			select {
			// goroutine ends when a timeout occurs
			case <-ctx.Context().Done():
				return
			default:
				for _, name := range defaultZKAuthOAuthProviders {
					provider := types.GetConfig(name)
					fetchIntervals = provider.FetchIntervals

					jwks, err := types.FetchJWK(provider.JwkEndpoint)
					if err != nil {
						time.Sleep(time.Duration(fetchIntervals) * time.Second)
						logger.Error(fmt.Sprintf("%s", err))
						continue
					}

					// add jwk
					for _, jwk := range jwks.Keys {
						jwkCopy := jwk
						k.jwksMap.AddJWK(&jwkCopy)
					}
				}
				time.Sleep(time.Duration(fetchIntervals) * time.Second)
			}
		}
	}()
}

func (k Keeper) GetJWKSize() int {
	return k.jwksMap.Size()
}

func (k Keeper) GetJWK(kid string) *types.JWK {
	return k.jwksMap.GetJWK(kid)
}

func (k Keeper) GetVerifier() *types.ZKAuthVerifier {
	return &k.zkVerifier
}
