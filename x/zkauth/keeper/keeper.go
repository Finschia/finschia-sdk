package keeper

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
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

// Currently, the FetchJWK interval is 60s, but this will change once the JWK rotation interval is known.
const fetchIntervals = time.Duration(60) * time.Second

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
	results := make([][]byte, 0, len(msgs))
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

		results = append(results, msgResp.Data)

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

func (k Keeper) FetchJWK(ctx sdk.Context, wg *sync.WaitGroup) {
	quit := make(chan struct{})
	ticker := time.NewTicker(fetchIntervals)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	defer wg.Done()

	go func() {
		<-sigCh
		close(quit)
	}()

	go func() {
		k.fetch(ctx)
		for {
			select {
			case <-ticker.C:
				ctx.Logger().Info(fmt.Sprintf("JWK fetch start in %s", types.ModuleName))
				k.fetch(ctx)
			case <-quit:
				ctx.Logger().Info(fmt.Sprintf("Received quite signal, fetch jwk exiting in %s", types.ModuleName))
				defer ticker.Stop()
				return
			// goroutine ends when a timeout occurs
			case <-ctx.Context().Done():
				return
			}
		}
	}()
}

func (k Keeper) fetch(ctx sdk.Context) {
	var defaultZKAuthOAuthProviders = [1]types.OidcProvider{types.Google}

	for _, name := range defaultZKAuthOAuthProviders {
		provider := types.GetConfig(name)
		jwks, err := types.FetchJWK(provider.JwkEndpoint)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("%s", err))
			continue
		}

		// add jwk
		for _, jwk := range jwks.Keys {
			jwkCopy := jwk
			k.jwksMap.AddJWK(&jwkCopy)
		}
	}
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
