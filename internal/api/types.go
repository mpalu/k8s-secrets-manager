package api

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type APIConfig struct {
	EnableMetrics bool
	EnableTracing bool
	EnableAuth    bool
	CorsEnabled   bool
	CorsOrigins   []string
}
