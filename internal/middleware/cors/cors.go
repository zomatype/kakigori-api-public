package cors

import (
	"fmt"
	"net/http"
	"strings"
)

const comma = ","

type Config struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
	MaxAge       int
}

// DefaultConfig デフォルトのCORS設定（ハッカソン用に寛容な設定）
func DefaultConfig() Config {
	return Config{
		AllowOrigins: []string{"*"}, // 本番環境では適切に制限すべき
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Content-Type",
		},
		MaxAge: 86400, // 24時間
	}
}

func Middleware(config Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(config.AllowOrigins) > 0 && config.AllowOrigins[0] == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			if r.Method == http.MethodOptions {
				if len(config.AllowMethods) > 0 {
					w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, comma))
				}
				if len(config.AllowHeaders) > 0 {
					w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, comma))
				}
				if config.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))
				}
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
