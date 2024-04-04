package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func BasicAuth(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println(".env 読み込みできませんでした")
			return
		}
		fmt.Println(".env 読み込み成功")

		id, pw, ok := r.BasicAuth()
		if !ok {
			return
		}

		envid := os.Getenv("BASIC_AUTH_USER_ID")
		envpw := os.Getenv("BASIC_AUTH_PASSWORD")

		if id != envid || pw != envpw || !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		fmt.Println("success")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
