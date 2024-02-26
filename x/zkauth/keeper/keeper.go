package keeper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
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
	jwks     []types.JWK
	lock     sync.RWMutex
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
) *Keeper {

	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		jwks:     nil,
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) LoopJWK(ctx sdk.Context, nodeHome string) {
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

					err := k.FetchJWK(provider.JwkEndpoint, nodeHome, name)
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

func (k *Keeper) FetchJWK(endpoint, nodeHome string, name types.OidcProvider) error {
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

	targetFile := filepath.Join(nodeHome, k.CreateJWKFileName(name))
	err = os.WriteFile(targetFile, dataBytes, 0o600)
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
