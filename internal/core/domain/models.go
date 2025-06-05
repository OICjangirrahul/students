package domain

import "time"

// 学生構造体：システムに登録されている学生を表現
type Student struct {
	// 学生の一意識別子
	ID int64 `json:"id,omitempty" swaggerignore:"true"`
	// 学生の氏名（必須）
	Name string `json:"name" binding:"required" example:"John Doe"`
	// 学生のメールアドレス（必須、メール形式）
	Email string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	// 学生の年齢（必須）
	Age int `json:"age" binding:"required" example:"25"`
	// 学生のパスワード（必須、最小6文字）
	Password string `json:"password,omitempty" binding:"required,min=6" example:"SecurePass123"`
	// アカウント作成日時
	CreatedAt time.Time `json:"created_at,omitempty" swaggerignore:"true"`
	// アカウント更新日時
	UpdatedAt time.Time `json:"updated_at,omitempty" swaggerignore:"true"`
}

// 教師構造体：システムに登録されている教師を表現
type Teacher struct {
	// 教師の一意識別子
	ID int64 `json:"id,omitempty" swaggerignore:"true"`
	// 教師の氏名（必須）
	Name string `json:"name" binding:"required" example:"Jane Smith"`
	// 教師のメールアドレス（必須、メール形式）
	Email string `json:"email" binding:"required,email" example:"jane.smith@example.com"`
	// 担当科目（必須）
	Subject string `json:"subject" binding:"required" example:"Mathematics"`
	// 教師のパスワード（必須、最小6文字）
	Password string `json:"password,omitempty" binding:"required,min=6" example:"SecurePass123"`
	// アカウント作成日時
	CreatedAt time.Time `json:"created_at,omitempty" swaggerignore:"true"`
	// アカウント更新日時
	UpdatedAt time.Time `json:"updated_at,omitempty" swaggerignore:"true"`
}

// 学生ログイン構造体：学生のログイン認証に使用
type StudentLogin struct {
	// 学生のメールアドレス（必須、メール形式）
	Email string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	// 学生のパスワード（必須、最小6文字）
	Password string `json:"password" binding:"required,min=6" example:"SecurePass123"`
}

// 教師ログイン構造体：教師のログイン認証に使用
type TeacherLogin struct {
	// 教師のメールアドレス（必須、メール形式）
	Email string `json:"email" binding:"required,email" example:"jane.smith@example.com"`
	// 教師のパスワード（必須、最小6文字）
	Password string `json:"password" binding:"required,min=6" example:"SecurePass123"`
}

// ファイルアップロードリクエスト構造体：S3へのファイルアップロード時に使用
type FileUpload struct {
	// ファイルの名前（必須）
	Name string `json:"name" binding:"required" example:"document.pdf"`
	// ファイルのバイナリデータ
	Data []byte `json:"data" swaggertype:"string" format:"binary"`
	// ファイルのMIMEタイプ
	ContentType string `json:"content_type" example:"application/pdf"`
}

// ファイル構造体：S3に保存されているファイルを表現
type File struct {
	// ファイルの一意識別子
	ID string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	// ファイルの名前
	Name string `json:"name" example:"document.pdf"`
	// ファイルのサイズ（バイト）
	Size int64 `json:"size" example:"1048576"`
	// ファイルのMIMEタイプ
	ContentType string `json:"content_type" example:"application/pdf"`
	// ファイルのURL
	URL string `json:"url" example:"https://storage.example.com/files/document.pdf"`
	// ファイルが保存されているS3バケット名
	BucketName string `json:"bucket_name" example:"my-bucket"`
	// ファイルのアップロード日時
	UploadedAt time.Time `json:"uploaded_at" example:"2024-03-21T15:30:45Z"`
}
