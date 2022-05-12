go 1.16

module github.com/line/lbm-sdk

require (
	github.com/99designs/keyring v1.1.6
	github.com/VictoriaMetrics/fastcache v1.10.0
	github.com/armon/go-metrics v0.3.11
	github.com/bgentry/speakeasy v0.1.0
	github.com/btcsuite/btcd v0.22.1
	github.com/coinbase/rosetta-sdk-go v0.7.9
	github.com/confio/ics23/go v0.7.0
	github.com/coocood/freecache v1.2.1
	github.com/cosmos/btcutil v1.0.4
	github.com/cosmos/go-bip39 v1.0.0
	github.com/cosmos/ledger-cosmos-go v0.11.1
	github.com/dgraph-io/ristretto v0.1.0
	github.com/dvsekhvalnov/jose2go v0.0.0-20200901110807-248326c1351b
	github.com/go-kit/kit v0.12.0
	github.com/gogo/gateway v1.1.0
	github.com/gogo/protobuf v1.3.3
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/gofuzz v1.1.1-0.20200604201612-c04b05f3adfa
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/hdevalence/ed25519consensus v0.0.0-20220222234857-c00d1f31bab3
	github.com/improbable-eng/grpc-web v0.15.0
	github.com/jhump/protoreflect v1.10.3
	github.com/line/iavl/v2 v2.0.0-init.1.0.20220215225951-cb11c91d8857
	github.com/line/ostracon v1.0.5
	github.com/line/tm-db/v2 v2.0.0-init.1.0.20220121012851-61d2bc1d9486
	github.com/line/wasmvm v0.16.3-0.9.0
	github.com/magiconair/properties v1.8.6
	github.com/mailru/easyjson v0.7.7
	github.com/mattn/go-isatty v0.0.14
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/prometheus/common v0.34.0
	github.com/rakyll/statik v0.1.7
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/rs/zerolog v1.26.1
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.11.0
	github.com/stretchr/testify v1.7.1
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	github.com/tendermint/btcd v0.1.1
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
	github.com/tendermint/go-amino v0.16.0
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4
	google.golang.org/genproto v0.0.0-20220407144326-9054f6ed7bac
	google.golang.org/grpc v1.46.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/99designs/keyring => github.com/cosmos/keyring v1.1.7-0.20210622111912-ef00f8ac3d76
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
