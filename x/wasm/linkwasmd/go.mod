module github.com/line/lfb-sdk/x/wasm/linkwasmd

go 1.13

require (
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d // indirect
	github.com/line/lfb-sdk v2.0.0-init.1.0.20210407071744-95eb4e7aef27
	github.com/line/ostracon v0.34.9-0.20210429084710-ef4fe0a40c7d
	github.com/line/tm-db/v2 v2.0.0-init.1.0.20210413083915-5bb60e117524
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20210305035536-64b5b1c73954 // indirect
	github.com/tendermint/tm-db v0.6.4 // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/line/lfb-sdk => ../../..
	github.com/tendermint/tm-db => github.com/line/tm-db v0.5.2
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
