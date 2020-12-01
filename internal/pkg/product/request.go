package product

type createRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Price       int    `json:"price"`
	Amount      int    `json:"amount"`
}
