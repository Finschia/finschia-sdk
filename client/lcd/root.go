package lcd

import (
	"time"

	"net/http"

	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	// unnamed import of statik for swagger UI support
	_ "github.com/link-chain/link/client/lcd/statik"
)

// to refresh cache per deployment
var eTag = time.Now().String()

// ServeCommand will start the application REST service as a blocking process. It
// takes a codec to create a RestServer object and a function to register all
// necessary routes.
func ServeCommand(cdc *codec.Codec, registerRoutesFn func(*lcd.RestServer)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rs := lcd.NewRestServer(cdc)

			registerRoutesFn(rs)
			registerSwaggerUI(rs)

			// Start the rest server and return error if one exists
			err = rs.Start(
				viper.GetString(flags.FlagListenAddr),
				viper.GetInt(flags.FlagMaxOpenConnections),
				uint(viper.GetInt(flags.FlagRPCReadTimeout)),
				uint(viper.GetInt(flags.FlagRPCWriteTimeout)),
			)

			return err
		},
	}

	return flags.RegisterRestServerFlags(cmd)
}

func registerSwaggerUI(rs *lcd.RestServer) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	staticServer := http.FileServer(statikFS)
	rs.Mux.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", cacheControlWrapper(staticServer)))
}

// no cache for the static server
func cacheControlWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=60") // cache allowed for 60 secs
		w.Header().Set("ETag", eTag)                  // to refresh cache per deployment
		h.ServeHTTP(w, r)
	})
}
