package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSミドルウェアを作成
// クロスオリジンリソース共有（CORS）の設定を行い、異なるドメインからのアクセスを制御する
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 全てのオリジンからのアクセスを許可
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 許可するHTTPメソッドを指定
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 許可するHTTPヘッダーを指定
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		// クライアントに公開するレスポンスヘッダーを指定
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Link")
		// クレデンシャル（認証情報）の送信を許可
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// プリフライトリクエストの結果をキャッシュする時間（秒）
		c.Writer.Header().Set("Access-Control-Max-Age", "300")

		// OPTIONSリクエスト（プリフライトリクエスト）の場合は204を返して終了
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
