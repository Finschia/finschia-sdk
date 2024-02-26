package keeper

import (
	"fmt"
	"time"

	"github.com/Finschia/finschia-sdk/codec"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
	"github.com/Finschia/ostracon/libs/log"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	jwksMap    *types.JWKsMap       // JWK manager
	zkVerifier types.ZKAuthVerifier // zkp verification key byte
	homePath   string               // root directory of app config
}

var _ types.ZKAuthKeeper = &Keeper{}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	jwksMap *types.JWKsMap,
	zkVerifier types.ZKAuthVerifier,
	homePath string,
) *Keeper {
	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		jwksMap:    jwksMap,
		zkVerifier: zkVerifier,
		homePath:   homePath,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
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
	return len(k.jwksMap.JWKs)
}

func (k Keeper) GetJWK(kid string) *types.JWK {
	return k.jwksMap.GetJWK(kid)
}

func (k Keeper) GetVerifier() *types.ZKAuthVerifier {
	return &k.zkVerifier
}
