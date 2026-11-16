package domain

type Rate struct {
	Currency string  `json:"currency"`
	Rate     float32 `json:"rate"`
}
