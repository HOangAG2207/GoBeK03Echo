package model

type ShortenURLRequest struct {
	URL string `json:"url" validate:"required,url" example:"http://localhost:8081/v1/docs/index.html"`
	Exp int64  `json:"exp" validate:"gt=0" example:"604800"` // thời gian hết hạn (giây)
}

type ShortenURLResponse struct {
	Code string `json:"code"`
}
