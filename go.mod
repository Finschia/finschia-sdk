module github.com/line/link-modules

go 1.13

require (
	github.com/CosmWasm/go-cosmwasm v0.10.0
	github.com/cosmos/cosmos-sdk v0.38.4
	github.com/golang/mock v1.4.3
	github.com/google/gofuzz v1.0.0
	github.com/gorilla/mux v1.7.4
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.4
	github.com/tendermint/tm-db v0.5.1
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/cosmos/cosmos-sdk => github.com/line/cosmos-sdk v0.38.4-0.1.0
