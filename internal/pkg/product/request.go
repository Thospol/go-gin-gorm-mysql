package product

type createRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Price       int    `json:"price" validate:"required"`
	Amount      int    `json:"amount" validate:"required"`
}
