package views

type UserRegistRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Street   string `json:"street"`
	City     string `json:"city"`
	Country  string `json:"country"`
	Kodepos  string `json:"kode pos"`
}
