package company

type UserRequest struct {
	Name     string `json:"name" validate:"required,alpha"`
	Password string `json:"password" validate:"required"`
}

type UserCreateResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UserLoginResponse struct {
	Hash string `json:"hash"`
}
