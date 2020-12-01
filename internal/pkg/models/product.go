package models

// Product product model
type Product struct {
	Model
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Price       int    `json:"price"`
	Amount      int    `json:"amount"`
}
