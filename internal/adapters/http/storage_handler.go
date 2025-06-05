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

type StorageHandler struct {
	fileStorage     ports.FileStorage
	documentStorage ports.DocumentStorage
	validator       *validator.Validate
}

func NewStorageHandler(fileStorage ports.FileStorage, documentStorage ports.DocumentStorage) *StorageHandler {
	return &StorageHandler{
		fileStorage:     fileStorage,
		documentStorage: documentStorage,
		validator:       validator.New(),
	}
}

// UploadFile handles file upload to S3
// @Summary      Upload a file to S3
// @Description  Upload a file to S3 storage
// @Tags         files
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "File to upload"
// @Success      201  {object}  domain.File
// @Router       /api/v1/files [post]
func (h *StorageHandler) UploadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid file")
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "failed to read file")
			return
		}

		uploadFile := &domain.FileUpload{
			Name:        header.Filename,
			ContentType: header.Header.Get("Content-Type"),
			Data:        data,
		}

		result, err := h.fileStorage.Upload(c.Request.Context(), uploadFile)
		if err != nil {
			slog.Error("failed to upload file", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to upload file")
			return
		}

		response.Success(c, http.StatusCreated, result)
	}
}

// DownloadFile handles file download from S3
// @Summary      Download a file from S3
// @Description  Download a file from S3 storage
// @Tags         files
// @Produce      octet-stream
// @Param        id path string true "File ID"
// @Success      200
// @Router       /api/v1/files/{id} [get]
func (h *StorageHandler) DownloadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		file, data, err := h.fileStorage.Download(c.Request.Context(), id)
		if err != nil {
			slog.Error("failed to download file", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to download file")
			return
		}

		c.Header("Content-Disposition", "attachment; filename="+file.Name)
		c.Data(http.StatusOK, file.ContentType, data)
	}
}

// ListFiles handles listing files from S3
// @Summary      List files from S3
// @Description  List all files from S3 storage
// @Tags         files
// @Produce      json
// @Success      200  {array}   domain.File
// @Router       /api/v1/files [get]
func (h *StorageHandler) ListFiles() gin.HandlerFunc {
	return func(c *gin.Context) {
		files, err := h.fileStorage.List(c.Request.Context())
		if err != nil {
			slog.Error("failed to list files", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to list files")
			return
		}

		response.Success(c, http.StatusOK, files)
	}
}

// DeleteFile handles file deletion from S3
// @Summary      Delete a file from S3
// @Description  Delete a file from S3 storage
// @Tags         files
// @Param        id path string true "File ID"
// @Success      204
// @Router       /api/v1/files/{id} [delete]
func (h *StorageHandler) DeleteFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := h.fileStorage.Delete(c.Request.Context(), id)
		if err != nil {
			slog.Error("failed to delete file", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to delete file")
			return
		}

		response.Success(c, http.StatusNoContent, nil)
	}
}

// CreateDocument handles document creation in DynamoDB
// @Summary      Create a document in DynamoDB
// @Description  Create a new document in DynamoDB
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        document body domain.DocumentCreate true "Document to create"
// @Success      201  {object}  domain.Document
// @Router       /api/v1/documents [post]
func (h *StorageHandler) CreateDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		var doc domain.DocumentCreate
		if err := c.ShouldBindJSON(&doc); err != nil {
			response.Error(c, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := h.validator.Struct(doc); err != nil {
			response.Error(c, http.StatusBadRequest, "validation failed")
			return
		}

		result, err := h.documentStorage.Create(c.Request.Context(), &doc)
		if err != nil {
			slog.Error("failed to create document", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to create document")
			return
		}

		response.Success(c, http.StatusCreated, result)
	}
}

// GetDocument handles getting a document from DynamoDB
// @Summary      Get a document from DynamoDB
// @Description  Get a document from DynamoDB by ID
// @Tags         documents
// @Produce      json
// @Param        id path string true "Document ID"
// @Success      200  {object}  domain.Document
// @Router       /api/v1/documents/{id} [get]
func (h *StorageHandler) GetDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		doc, err := h.documentStorage.Get(c.Request.Context(), id)
		if err != nil {
			slog.Error("failed to get document", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to get document")
			return
		}

		response.Success(c, http.StatusOK, doc)
	}
}

// UpdateDocument handles document update in DynamoDB
// @Summary      Update a document in DynamoDB
// @Description  Update a document in DynamoDB by ID
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        id path string true "Document ID"
// @Param        document body domain.DocumentUpdate true "Document update data"
// @Success      200  {object}  domain.Document
// @Router       /api/v1/documents/{id} [put]
func (h *StorageHandler) UpdateDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var update domain.DocumentUpdate
		if err := c.ShouldBindJSON(&update); err != nil {
			response.Error(c, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := h.validator.Struct(update); err != nil {
			response.Error(c, http.StatusBadRequest, "validation failed")
			return
		}

		result, err := h.documentStorage.Update(c.Request.Context(), id, &update)
		if err != nil {
			slog.Error("failed to update document", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to update document")
			return
		}

		response.Success(c, http.StatusOK, result)
	}
}

// DeleteDocument handles document deletion from DynamoDB
// @Summary      Delete a document from DynamoDB
// @Description  Delete a document from DynamoDB by ID
// @Tags         documents
// @Param        id path string true "Document ID"
// @Success      204
// @Router       /api/v1/documents/{id} [delete]
func (h *StorageHandler) DeleteDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := h.documentStorage.Delete(c.Request.Context(), id)
		if err != nil {
			slog.Error("failed to delete document", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to delete document")
			return
		}

		response.Success(c, http.StatusNoContent, nil)
	}
}

// ListDocuments handles listing documents from DynamoDB
// @Summary      List documents from DynamoDB
// @Description  List documents from DynamoDB by type
// @Tags         documents
// @Produce      json
// @Param        type query string true "Document type"
// @Success      200  {array}   domain.Document
// @Router       /api/v1/documents [get]
func (h *StorageHandler) ListDocuments() gin.HandlerFunc {
	return func(c *gin.Context) {
		docType := c.Query("type")
		if docType == "" {
			response.Error(c, http.StatusBadRequest, "document type is required")
			return
		}

		docs, err := h.documentStorage.List(c.Request.Context(), docType)
		if err != nil {
			slog.Error("failed to list documents", slog.String("error", err.Error()))
			response.Error(c, http.StatusInternalServerError, "failed to list documents")
			return
		}

		response.Success(c, http.StatusOK, docs)
	}
}
