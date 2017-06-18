package api

type CreateContactResponse struct {
	Id int `json:"id"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type EmptyObject struct {
}
