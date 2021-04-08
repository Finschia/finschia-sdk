module "github.com/line/lbm-sdk/v2/x/wasm/linkwasmd"

go 1.13

require (
	github.com/cosmos/cosmos-sdk v0.39.2
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/line/lbm-sdk/v2/ v0.2.0
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/rakyll/statik v0.1.7 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.9
	github.com/tendermint/tm-db v0.5.2
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

replace (
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/tendermint/tendermint => github.com/line/linemint v1.0.0
	github.com/tendermint/tm-db => github.com/line/tm-db v0.5.2
	github.com/cosmos/cosmos-sdk => github.com/line/lbm-sdk v0.39.2-0.2.0
	"github.com/line/lbm-sdk/v2/ => ../../..
)