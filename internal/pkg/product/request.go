package product

type createRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Price       int    `json:"price" validate:"required"`
	Amount      int    `json:"amount" validate:"required"`
}

type updateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Amount      int    `json:"amount"`
}
