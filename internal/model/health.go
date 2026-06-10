package model

type HealthCheckResponse struct {
	Message     string `json:"message" example:"OK"`
	ServiceName string `json:"service_name" example:"GoBe K03 API"`
	InstanceID  string `json:"instance_id" example:"instance-12345"`
}
type HealthCheckSwaggerResponse struct {
	Status      string `json:"status"`
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id"`
}
