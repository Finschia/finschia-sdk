package loadgenerator

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/line/link/contrib/load_test/types"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func RegisterHandlers(lg *LoadGenerator, r *mux.Router) {
	r.HandleFunc("/target/load", LoadHandler(lg)).Methods("POST")
	r.HandleFunc("/target/fire", FireHandler(lg)).Methods("POST")
}

func LoadHandler(lg *LoadGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoadRequest

		if err := parseRequest(r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Config.TPS <= 0 || req.Config.Duration <= 0 || req.Config.TargetURL == "" || req.Config.ChainID == "" {
			http.Error(w, types.InvalidLoadParameterError.Error("invalid parameter of load handler"), http.StatusBadRequest)
			return
		}
		if req.Config.PacerType != types.ConstantPacer && req.Config.PacerType != types.LinearPacer {
			http.Error(w, types.InvalidPacerTypeError{PacerType: req.Config.PacerType}.Error(), http.StatusBadRequest)
			return
		}

		if err := lg.ApplyConfig(req.Config); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch req.TargetType {
		case types.QueryAccount:
			if err := lg.RunWithGoroutines(lg.GenerateAccountQueryTarget); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case types.TxSend:
			if err := lg.RunWithGoroutines(lg.GenerateMsgSendTxTarget); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, types.InvalidTargetTypeError.Error("invalid target type"), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func parseRequest(r *http.Request, req interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &req)
	if err != nil {
		return err
	}
	return nil
}

func FireHandler(lg *LoadGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enc := vegeta.NewEncoder(w)
		for res := range lg.Fire(r.Host) {
			if err := enc.Encode(res); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
