package keeper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Finschia/finschia-sdk/codec"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
	"github.com/Finschia/ostracon/libs/log"
)

const JwkFileName = "jwk.json"

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

func (k Keeper) FetchJWK(ctx sdk.Context, nodeHome string) {
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

					resp, err := k.GetJWK(provider.JwkEndpoint)
					if err != nil {
						time.Sleep(time.Duration(fetchIntervals) * time.Second)
						logger.Error(fmt.Sprintf("%s", err))
						continue
					}

					file, err := os.Create(filepath.Join(nodeHome, k.CreateJWKFileName(name)))
					if err != nil {
						time.Sleep(time.Duration(fetchIntervals) * time.Second)
						logger.Error(fmt.Sprintf("%s", err))
						continue
					}

					_, err = io.Copy(file, resp.Body)
					if err != nil {
						time.Sleep(time.Duration(fetchIntervals) * time.Second)
						logger.Error(fmt.Sprintf("%s", err))
						continue
					}

					resp.Body.Close()
					file.Close()
				}
				time.Sleep(time.Duration(fetchIntervals) * time.Second)
			}
		}
	}()
}

func (k Keeper) ParseJWKs(byteArray []byte) (jwks []types.JWK, err error) {
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

func (k Keeper) GetJWK(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (k Keeper) CreateJWKFileName(name types.OidcProvider) string {
	fileNamePattern := strings.Replace(JwkFileName, ".", "-%s.", 1)
	return fmt.Sprintf(fileNamePattern, name)
}
