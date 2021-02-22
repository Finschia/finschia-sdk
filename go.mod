module github.com/line/link-modules

go 1.13

require (
	github.com/CosmWasm/wasmvm v0.12.0
	github.com/cosmos/cosmos-sdk v0.39.2
	github.com/golang/mock v1.4.3
	github.com/google/gofuzz v1.0.0
	github.com/gorilla/mux v1.7.4
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.9
	github.com/tendermint/tm-db v0.5.2
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/CosmWasm/wasmvm => github.com/line/wasmvm v0.12.0-0.1.0

replace github.com/cosmos/cosmos-sdk => github.com/line/cosmos-sdk v0.39.2-0.2.0
