package types

type LoadRequest struct {
	TargetType string `json:"target_type"`
	Config     Config `json:"config"`
}

func NewLoadRequest(targetType string, config Config) *LoadRequest {
	return &LoadRequest{
		TargetType: targetType,
		Config:     config,
	}
}
