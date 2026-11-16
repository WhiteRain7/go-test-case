package domain

type UUID = string

type Wallet struct {
	ID   UUID `json:"id"`
	User int  `json:"user"`
}

type WalletCurrency struct {
	Wallet   UUID    `json:"wallet"`
	Currency string  `json:"currency"`
	Amount   float32 `json:"amount"`
}

type WalletBalance struct {
	Wallet  UUID
	Balance map[string]float32 `json:"balance"`
}
