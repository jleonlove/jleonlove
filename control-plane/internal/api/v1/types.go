package v1

type Action struct {
	Tool string `json:"tool"`
	Operation string `json:"operation"`
	Arguments map[string]any `json:"arguments"`
}

type ExecuteRequest struct {
	Envelope []byte `json:"envelope"`
	Action Action `json:"action"`
}

type ExecuteResponse struct {
	ExecutionID string `json:"execution_id"`
	Status string `json:"status"`
	Result any `json:"result,omitempty"`
}
