package domain

type Student struct {
	ID       int64  `json:"id"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Age      int    `json:"age" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type StudentLogin struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
