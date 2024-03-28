package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mileusna/useragent"
)

type key string

var OSkey key = "os"

func OS(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("os middleware")

		ua := useragent.Parse(r.UserAgent())                // httpリクエストからユーザー情報取得
		ctx := context.WithValue(r.Context(), OSkey, ua.OS) // コンテキストに値を格納,keyは専用の型を定義するかつ変数の必要あり

		fmt.Println(ctx.Value(OSkey))      // OS情報を出力(Windows)
		h.ServeHTTP(w, r.WithContext(ctx)) // httpリクエストにコンテキスト,OS情報をセット
	}
	return http.HandlerFunc(fn)
}
