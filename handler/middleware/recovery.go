package middleware

import (
	"fmt"
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: ここに実装をする
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}() // 関数を実行(deferなので後から)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
