module github.com/line/link

go 1.13

require (
	github.com/blend/go-sdk v2.0.0+incompatible // indirect
	github.com/cosmos/cosmos-sdk v0.38.3
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/mock v1.4.3
	github.com/gonum/stat v0.0.0-20181125101827-41a0da705a5b
	github.com/gorilla/mux v1.7.4
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/olekukonko/tablewriter v0.0.4
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/otiai10/copy v1.1.1
	github.com/rakyll/statik v0.1.7
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0
	github.com/snikch/goodman v0.0.0-20171125024755-10e37e294daa
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.4
	github.com/tendermint/tm-db v0.5.1
	github.com/tsenart/vegeta/v12 v12.8.3
	github.com/wcharczuk/go-chart v2.0.1+incompatible
	golang.org/x/image v0.0.0-20200119044424-58c23975cae1 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v3 v3.0.0-20191120175047-4206685974f2
)

replace github.com/cosmos/cosmos-sdk => github.com/line/cosmos-sdk v0.38.3-t0.33.4-0.0.0
