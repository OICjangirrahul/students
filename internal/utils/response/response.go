package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// レスポンス構造体：APIレスポンスの基本構造を定義
type Response struct {
	// ステータス：処理の結果を示す（"OK" または "Error"）
	Status string `json:"status"`
	// エラー：エラーメッセージを格納（エラーがある場合のみ）
	Error string `json:"error"`
}

// レスポンスステータスの定数
const (
	// 正常終了を示すステータス
	StatusOk = "OK"
	// エラー発生を示すステータス
	StatusError = "Error"
)

// JSONレスポンスを書き込む
// HTTPレスポンスライターにJSONデータを書き込み、ステータスコードを設定する
func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

// 一般的なエラーレスポンスを生成
// エラーオブジェクトからエラーメッセージを抽出してレスポンスを作成
func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

// バリデーションエラーレスポンスを生成
// バリデーションエラーの詳細を解析し、人間が読みやすいエラーメッセージを作成
func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ","),
	}
}

// 成功レスポンスを送信
// 指定されたステータスコードとデータでJSONレスポンスを生成
func Success(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{
		"success": true,
		"data":    data,
	})
}

// エラーレスポンスを送信
// 指定されたステータスコードとエラーメッセージでJSONレスポンスを生成
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"success": false,
		"error":   message,
	})
}
