package keeper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Finschia/finschia-sdk/codec"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
	"github.com/Finschia/ostracon/libs/log"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
) *Keeper {

	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetJWK(ctx sdk.Context, kid string) types.JWK {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(kid))
	var jwk types.JWK
	if bz == nil {
		return jwk
	}

	k.cdc.MustUnmarshal(bz, &jwk)
	return jwk
}

func (k Keeper) SetJWK(ctx sdk.Context, kid string, jwk types.JWK) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&jwk)
	if err != nil {
		return err
	}

	store.Set([]byte(kid), bz)
	return nil
}

func (k Keeper) GetKidList(ctx sdk.Context, iss string) types.JwkIdList {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(iss))
	var kidList types.JwkIdList
	if bz == nil {
		return kidList
	}

	k.cdc.MustUnmarshal(bz, &kidList)
	return kidList
}

func (k Keeper) SetKidList(ctx sdk.Context, kidList types.JwkIdList) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&kidList)
	if err != nil {
		return err
	}

	store.Set([]byte(kidList.Iss), bz)
	return nil
}

func (k Keeper) FetchJWK(ctx sdk.Context) {
	logger := k.Logger(ctx)
	var defaultZKAuthOAuthProviders = [1]types.OidcProvider{types.Google}
	var fetchIntervals uint64
	go func() {
		for {
			for _, name := range defaultZKAuthOAuthProviders {
				provider := types.GetConfig(name)
				fetchIntervals = provider.FetchIntervals
				req, err := http.NewRequest("GET", provider.JwkEndpoint, nil)
				if err != nil {
					logger.Error(fmt.Sprintf("%s", err))
				}
				client := new(http.Client)
				resp, err := client.Do(req)
				if err != nil {
					logger.Error(fmt.Sprintf("%s", err))
				}

				resp.Body.Close()

				byteArray, _ := io.ReadAll(resp.Body)

				var data map[string]interface{}
				json.Unmarshal(byteArray, &data)

				for _, v := range data["keys"].([]interface{}) {
					var jwk types.JWK

					kid := v.(map[string]interface{})["kid"].(string)

					b, err := json.Marshal(v)
					if err != nil {
						logger.Error(fmt.Sprintf("%s", err))
					}

					if err := json.Unmarshal(b, &jwk); err != nil {
						logger.Error(fmt.Sprintf("%s", err))
					}

					var kidList types.JwkIdList
					var kids []string
					kids = append(kids, kid)
					kidList.Iss = provider.Iss
					kidList.Kid = kids

					// Set jwk info
					k.SetKidList(ctx, kidList)
					k.SetJWK(ctx, kid, jwk)
				}
			}
			time.Sleep(time.Duration(fetchIntervals) * time.Second)
		}
	}()
}
