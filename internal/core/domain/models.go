package domain

import "time"

// Student represents a student in the system
type Student struct {
	ID        int64     `json:"id,omitempty" swaggerignore:"true"`
	Name      string    `json:"name" binding:"required" example:"John Doe"`
	Email     string    `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Age       int       `json:"age" binding:"required" example:"25"`
	Password  string    `json:"password,omitempty" binding:"required,min=6" example:"SecurePass123"`
	CreatedAt time.Time `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt time.Time `json:"updated_at,omitempty" swaggerignore:"true"`
}

// Teacher represents a teacher in the system
type Teacher struct {
	ID        int64     `json:"id,omitempty" swaggerignore:"true"`
	Name      string    `json:"name" binding:"required" example:"Jane Smith"`
	Email     string    `json:"email" binding:"required,email" example:"jane.smith@example.com"`
	Subject   string    `json:"subject" binding:"required" example:"Mathematics"`
	Password  string    `json:"password,omitempty" binding:"required,min=6" example:"SecurePass123"`
	CreatedAt time.Time `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt time.Time `json:"updated_at,omitempty" swaggerignore:"true"`
}

// StudentLogin represents student login credentials
type StudentLogin struct {
	Email    string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"SecurePass123"`
}

// TeacherLogin represents teacher login credentials
type TeacherLogin struct {
	Email    string `json:"email" binding:"required,email" example:"jane.smith@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"SecurePass123"`
}

// FileUpload represents a file upload request
type FileUpload struct {
	Name        string `json:"name" binding:"required" example:"document.pdf"`
	Data        []byte `json:"data" swaggertype:"string" format:"binary"`
	ContentType string `json:"content_type" example:"application/pdf"`
}

// File represents a stored file
type File struct {
	ID          string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string    `json:"name" example:"document.pdf"`
	Size        int64     `json:"size" example:"1048576"`
	ContentType string    `json:"content_type" example:"application/pdf"`
	URL         string    `json:"url" example:"https://storage.example.com/files/document.pdf"`
	BucketName  string    `json:"bucket_name" example:"my-bucket"`
	UploadedAt  time.Time `json:"uploaded_at" example:"2024-03-21T15:30:45Z"`
}
