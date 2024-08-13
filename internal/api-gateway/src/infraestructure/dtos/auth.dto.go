package dtos

type (
	SignUpDto struct {
		Password string `json:"password" validate:"required,min=8,max=20"`
		Name     string `json:"name" validate:"required,min=3,max=50"`
		Email    string `json:"email" validate:"required,email"`
		LasName  string `json:"lastName" validate:"required,min=3,max=50"`

		Gender string `json:"gender"`
	}

	SignInDto struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=20"`
	}
)
