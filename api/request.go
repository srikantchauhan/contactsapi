package api

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateContactRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type AddPhoneRequest struct {
	Phone string `json:"phone"`
}
