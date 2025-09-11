package requests

type RegisterRequest struct {
    Username string `json:"username" validate:"required,min=3"`
    Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
    Username string `json:"username" validate:"required,min=3"`
    Password string `json:"password" validate:"required,min=6"`
}

func (r *RegisterRequest) Validate() error {
	return ValidateStruct(r)
}

func (r *LoginRequest) Validate() error {
	return ValidateStruct(r)
}