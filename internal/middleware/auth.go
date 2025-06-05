package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/OICjangirrahul/students/internal/config"
	"github.com/OICjangirrahul/students/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// ロールベースのアクセス制御ミドルウェアを作成
// 指定された役割（ロール）を持つユーザーのみがアクセスを許可される
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// コンテキストから役割を取得
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("role not found in token")))
			c.Abort()
			return
		}

		// 役割を文字列に変換
		roleStr, ok := role.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid role type")))
			c.Abort()
			return
		}

		// 許可された役割かどうかを確認
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, response.GeneralError(fmt.Errorf("access denied: insufficient privileges")))
		c.Abort()
	}
}

// JWT認証ミドルウェアを作成
// リクエストヘッダーからJWTトークンを検証し、ユーザー情報をコンテキストに追加
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 認証ヘッダーを取得
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("authorization header is required")))
			c.Abort()
			return
		}

		// トークン文字列を抽出
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		// トークンを検証
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid token: %v", err)))
			c.Abort()
			return
		}

		// トークンの有効性とクレームを確認
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// クレームをコンテキストに追加
			c.Set("userID", claims["sub"])
			c.Set("email", claims["email"])
			c.Set("role", claims["role"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid token claims")))
			c.Abort()
			return
		}
	}
}

// リソース所有権チェックミドルウェアを作成
// ユーザーが自分のリソースにのみアクセスできるようにする
// ただし、教師は全てのリソースにアクセス可能
func ResourceOwnershipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// コンテキストからユーザーIDを取得
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("user ID not found in token")))
			c.Abort()
			return
		}

		// URLパラメータからリソースIDを取得
		resourceID := c.Param("id")
		if resourceID == "" {
			c.Next()
			return
		}

		// ユーザーIDを文字列に変換して比較
		userIDStr := fmt.Sprint(userID)

		// ユーザーが自分のリソースにアクセスしていない場合
		if userIDStr != resourceID {
			role, _ := c.Get("role")
			// 教師は全てのリソースにアクセス可能
			if roleStr, ok := role.(string); ok && roleStr == "teacher" {
				c.Next()
				return
			}
			c.JSON(http.StatusForbidden, response.GeneralError(fmt.Errorf("access denied: you can only access your own resources")))
			c.Abort()
			return
		}

		c.Next()
	}
}
