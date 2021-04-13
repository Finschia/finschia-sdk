module "github.com/line/lbm-sdk/v2/x/wasm/linkwasmd"

go 1.13

require (
	github.com/line/lbm-sdk/v2 v2.0.0-init.1.0.20210407071744-95eb4e7aef27
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/rakyll/statik v0.1.7 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/go-amino v0.15.1
	github.com/line/ostracon v0.34.9-0.20210406083837-4183d649b30c
	github.com/line/tm-db/v2 v2.0.0-init.1.0.20210406062110-9424ca70955a
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
	github.com/line/lbm-sdk/v2/ => ../../..
)
