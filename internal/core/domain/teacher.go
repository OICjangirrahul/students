package domain

type Teacher struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password,omitempty" validate:"required"`
	Subject  string    `json:"subject" validate:"required"`
	Students []Student `json:"students,omitempty"`
}

type TeacherLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
