package model

type HealthCheckResponse struct {
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id"`
}

// Swagger success response (flatten)
type HealthCheckSwaggerResponse struct {
	Status string `json:"status"`
	Info   string `json:"info"`
	HealthCheckResponse
}
