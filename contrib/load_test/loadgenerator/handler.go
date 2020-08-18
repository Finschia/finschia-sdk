package loadgenerator

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/line/link/contrib/load_test/scenario"
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

		if req.Config.TPS <= 0 || req.Config.Duration <= 0 || req.Config.TargetURL == "" ||
			req.Config.ChainID == "" || req.Config.RampUpTime < 0 {
			http.Error(w, types.InvalidLoadParameterError.Error("invalid parameter of load handler"), http.StatusBadRequest)
			return
		}

		if err := lg.ApplyConfig(req.Config); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		testScenario, ok := scenario.NewScenarios(req.Config, req.StateParams, req.ScenarioParams)[req.Scenario]
		if !ok {
			http.Error(w, types.InvalidScenarioError.Error("invalid scenario"), http.StatusBadRequest)
			return
		}
		if err := lg.RunWithGoroutines(testScenario.GenerateTarget); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
		var results []vegeta.Result
		for res := range lg.Fire(r.Host) {
			results = append(results, *res)
		}
		data, err := json.Marshal(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(data)
		if err != nil {
			log.Println("Failed to write results:", err)
			return
		}
	}
}
