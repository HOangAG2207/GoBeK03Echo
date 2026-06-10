package model

type ShortenURLRequest struct {
	URL string `json:"url" validate:"required,url" example:"https://google.com"`
	Exp int64  `json:"exp" validate:"gt=0" example:"604800"` // thời gian hết hạn (giây)
}

type ShortenURLResponse struct {
	Code string `json:"code"`
}
type ShortenURLSwaggerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
}
