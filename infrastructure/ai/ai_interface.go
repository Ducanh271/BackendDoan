package ai

type AIRegisterResponse struct {
	Status    string    `json:"status"`
	Embedding []float64 `json:"embedding"`
	Error     string    `json:"error"`
	Stats     struct {
		TotalReceived int      `json:"total_received"`
		ValidFaces    int      `json:"valid_faces"`
		LivenessPass  int      `json:"liveness_pass"`
		DeepfakePass  int      `json:"deepfake_pass"`
		Errors        []string `json:"errors"`
	} `json:"stats"`
}

type AIClient interface {
	RegisterFace(images [][]byte) (*AIRegisterResponse, error)
}
