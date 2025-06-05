package http

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/OICjangirrahul/students/internal/core/domain"
	"github.com/OICjangirrahul/students/internal/core/ports"
	"github.com/OICjangirrahul/students/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ストレージハンドラー構造体：ファイルとドキュメントの操作を管理
type StorageHandler struct {
	// ファイルストレージインターフェース
	fileStorage ports.FileStorage
	// ドキュメントストレージインターフェース
	documentStorage ports.DocumentStorage
	// バリデーター
	validator *validator.Validate
}

// 新しいストレージハンドラーを作成する関数
func NewStorageHandler(fileStorage ports.FileStorage, documentStorage ports.DocumentStorage) *StorageHandler {
	return &StorageHandler{
		fileStorage:     fileStorage,
		documentStorage: documentStorage,
		validator:       validator.New(),
	}
}

// S3にファイルをアップロードする機能を提供するハンドラー
// ファイルを受け取り、S3に保存し、保存されたファイルの情報を返す
// @Summary      Upload a file to S3
// @Description  Upload a file to S3 storage
// @Tags         files
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        file formData file true "File to upload"
// @Success      201  {object}  domain.File
// @Failure      401  {object}  response.Response "Unauthorized"
// @Failure      403  {object}  response.Response "Forbidden - Teacher role required"
// @Router       /api/v1/files [post]
func (h *StorageHandler) UploadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// マルチパートフォームからファイルを取得
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid file")
			return
		}
		defer file.Close()

		// ファイルデータを読み込む
		data, err := io.ReadAll(file)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "failed to read file")
			return
		}

		// アップロード用のファイル構造体を作成
		uploadFile := &domain.FileUpload{
			Name:        header.Filename,
			ContentType: header.Header.Get("Content-Type"),
			Data:        data,
		}

		// ファイルをS3にアップロード
		result, err := h.fileStorage.Upload(c.Request.Context(), uploadFile)
		if err != nil {
			slog.Error("failed to upload file", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to upload file")
			return
		}

		response.Success(c, http.StatusCreated, result)
	}
}

// S3からファイルをダウンロードする機能を提供するハンドラー
// ファイルIDを受け取り、対応するファイルをダウンロードする
// @Summary      Download a file from S3
// @Description  Download a file from S3 storage
// @Tags         files
// @Produce      octet-stream
// @Security     BearerAuth
// @Param        id path string true "File ID"
// @Success      200
// @Failure      401  {object}  response.Response "Unauthorized"
// @Failure      403  {object}  response.Response "Forbidden - Teacher role required"
// @Router       /api/v1/files/{id} [get]
func (h *StorageHandler) DownloadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// パスパラメータからファイルIDを取得
		id := c.Param("id")
		// S3からファイルをダウンロード
		file, data, err := h.fileStorage.Download(c.Request.Context(), id)
		if err != nil {
			slog.Error("failed to download file", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to download file")
			return
		}

		// ファイルをクライアントに送信
		c.Header("Content-Disposition", "attachment; filename="+file.Name)
		c.Data(http.StatusOK, file.ContentType, data)
	}
}

// S3に保存されているファイルの一覧を取得する機能を提供するハンドラー
// @Summary      List files from S3
// @Description  List all files from S3 storage
// @Tags         files
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   domain.File
// @Failure      401  {object}  response.Response "Unauthorized"
// @Failure      403  {object}  response.Response "Forbidden - Teacher role required"
// @Router       /api/v1/files [get]
func (h *StorageHandler) ListFiles() gin.HandlerFunc {
	return func(c *gin.Context) {
		// S3からファイル一覧を取得
		files, err := h.fileStorage.List(c.Request.Context())
		if err != nil {
			slog.Error("failed to list files", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to list files")
			return
		}

		response.Success(c, http.StatusOK, files)
	}
}

// S3からファイルを削除する機能を提供するハンドラー
// ファイルIDを受け取り、対応するファイルを削除する
// @Summary      Delete a file from S3
// @Description  Delete a file from S3 storage
// @Tags         files
// @Security     BearerAuth
// @Param        id path string true "File ID"
// @Success      204
// @Failure      401  {object}  response.Response "Unauthorized"
// @Failure      403  {object}  response.Response "Forbidden - Teacher role required"
// @Router       /api/v1/files/{id} [delete]
func (h *StorageHandler) DeleteFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// パスパラメータからファイルIDを取得
		id := c.Param("id")
		// S3からファイルを削除
		err := h.fileStorage.Delete(c.Request.Context(), id)
		if err != nil {
			slog.Error("failed to delete file", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to delete file")
			return
		}

		response.Success(c, http.StatusNoContent, nil)
	}
}

// DynamoDBに新しいドキュメントを作成する機能を提供するハンドラー
// ドキュメントデータを受け取り、新しいドキュメントを作成する
// @Summary      Create a document in DynamoDB
// @Description  Create a new document in DynamoDB
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        document body domain.DocumentCreate true "Document to create"
// @Success      201  {object}  domain.Document
// @Failure      401  {object}  response.Response "Unauthorized"
// @Failure      403  {object}  response.Response "Forbidden - Teacher role required"
// @Router       /api/v1/documents [post]
func (h *StorageHandler) CreateDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエストボディからドキュメントデータを取得
		var doc domain.DocumentCreate
		if err := c.ShouldBindJSON(&doc); err != nil {
			response.Error(c, http.StatusBadRequest, "invalid request body")
			return
		}

		// バリデーション実行
		if err := h.validator.Struct(doc); err != nil {
			response.Error(c, http.StatusBadRequest, "validation failed")
			return
		}

		// DynamoDBにドキュメントを作成
		result, err := h.documentStorage.Create(c.Request.Context(), &doc)
		if err != nil {
			slog.Error("failed to create document", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to create document")
			return
		}

		response.Success(c, http.StatusCreated, result)
	}
}

// DynamoDBからドキュメントを取得する機能を提供するハンドラー
// ドキュメントIDを受け取り、対応するドキュメントを返す
// @Summary      Get a document from DynamoDB
// @Description  Get a document from DynamoDB by ID
// @Tags         documents
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Document ID"
// @Success      200  {object}  domain.Document
// @Failure      401  {object}  response.Response "Unauthorized"
// @Failure      403  {object}  response.Response "Forbidden - Teacher role required"
// @Router       /api/v1/documents/{id} [get]
func (h *StorageHandler) GetDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		// パスパラメータからドキュメントIDを取得
		id := c.Param("id")
		// DynamoDBからドキュメントを取得
		doc, err := h.documentStorage.Get(c.Request.Context(), id)
		if err != nil {
			slog.Error("failed to get document", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to get document")
			return
		}

		response.Success(c, http.StatusOK, doc)
	}
}

// DynamoDBのドキュメントを更新する機能を提供するハンドラー
// ドキュメントIDと更新データを受け取り、ドキュメントを更新する
// @Summary      Update a document in DynamoDB
// @Description  Update a document in DynamoDB by ID
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Document ID"
// @Param        document body domain.DocumentUpdate true "Document update data"
// @Success      200  {object}  domain.Document
// @Failure      401  {object}  response.Response "Unauthorized"
// @Failure      403  {object}  response.Response "Forbidden - Teacher role required"
// @Router       /api/v1/documents/{id} [put]
func (h *StorageHandler) UpdateDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		// パスパラメータからドキュメントIDを取得
		id := c.Param("id")
		// リクエストボディから更新データを取得
		var update domain.DocumentUpdate
		if err := c.ShouldBindJSON(&update); err != nil {
			response.Error(c, http.StatusBadRequest, "invalid request body")
			return
		}

		// バリデーション実行
		if err := h.validator.Struct(update); err != nil {
			response.Error(c, http.StatusBadRequest, "validation failed")
			return
		}

		// DynamoDBのドキュメントを更新
		result, err := h.documentStorage.Update(c.Request.Context(), id, &update)
		if err != nil {
			slog.Error("failed to update document", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to update document")
			return
		}

		response.Success(c, http.StatusOK, result)
	}
}

// DynamoDBからドキュメントを削除する機能を提供するハンドラー
// ドキュメントIDを受け取り、対応するドキュメントを削除する
// @Summary      Delete a document from DynamoDB
// @Description  Delete a document from DynamoDB by ID
// @Tags         documents
// @Security     BearerAuth
// @Param        id path string true "Document ID"
// @Success      204
// @Failure      401  {object}  response.Response "Unauthorized"
// @Failure      403  {object}  response.Response "Forbidden - Teacher role required"
// @Router       /api/v1/documents/{id} [delete]
func (h *StorageHandler) DeleteDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		// パスパラメータからドキュメントIDを取得
		id := c.Param("id")
		// DynamoDBからドキュメントを削除
		err := h.documentStorage.Delete(c.Request.Context(), id)
		if err != nil {
			slog.Error("failed to delete document", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to delete document")
			return
		}

		response.Success(c, http.StatusNoContent, nil)
	}
}

// DynamoDBに保存されているドキュメントの一覧を取得する機能を提供するハンドラー
// ドキュメントタイプを受け取り、該当するドキュメントの一覧を返す
// @Summary      List documents from DynamoDB
// @Description  List documents from DynamoDB by type
// @Tags         documents
// @Produce      json
// @Security     BearerAuth
// @Param        type query string true "Document type"
// @Success      200  {array}   domain.Document
// @Failure      401  {object}  response.Response "Unauthorized"
// @Failure      403  {object}  response.Response "Forbidden - Teacher role required"
// @Router       /api/v1/documents [get]
func (h *StorageHandler) ListDocuments() gin.HandlerFunc {
	return func(c *gin.Context) {
		// クエリパラメータからドキュメントタイプを取得
		docType := c.Query("type")
		if docType == "" {
			response.Error(c, http.StatusBadRequest, "document type is required")
			return
		}

		// DynamoDBからドキュメント一覧を取得
		docs, err := h.documentStorage.List(c.Request.Context(), docType)
		if err != nil {
			slog.Error("failed to list documents", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to list documents")
			return
		}

		response.Success(c, http.StatusOK, docs)
	}
}
