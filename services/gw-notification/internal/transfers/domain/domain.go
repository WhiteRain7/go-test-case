package domain

type Transfer struct {
	User      int     `json:"user"`
	Operation string  `json:"operation"`
	Amount    float32 `json:"amount"`
	Currency  string  `json:"currency"`
	Timestamp string  `json:"timestamp"`
}
