package types

type LoadRequest struct {
	Scenario       string            `json:"scenario"`
	ScenarioParams []string          `json:"scenario_params"`
	Config         Config            `json:"config"`
	StateParams    map[string]string `json:"state_params"`
}

func NewLoadRequest(scenario string, scenarioParams []string, config Config, params map[string]string) *LoadRequest {
	return &LoadRequest{
		Scenario:       scenario,
		ScenarioParams: scenarioParams,
		Config:         config,
		StateParams:    params,
	}
}
