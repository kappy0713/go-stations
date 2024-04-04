package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/mileusna/useragent"
)

func Log(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("log middleware")

		accessTime := time.Now() // 現在時刻を取得
		urlPath := r.URL.Path    // Path取得

		ua := useragent.Parse(r.UserAgent())
		ctx := context.WithValue(r.Context(), OSkey, ua.OS)

		h.ServeHTTP(w, r)

		defer func() {
			fmt.Println("log defer")

			// OS情報取得
			os, ok := r.WithContext(ctx).Context().Value(OSkey).(string) // 同じmiddlewareだからここでOSkeyを定義する必要はない,stringにキャストする(Valueの型がanyだから)
			if !ok {
				log.Println("key does not exist")
				return
			}

			// 処理時間の計算,Millsecondsに変換
			end := time.Now()
			latency := int64(math.Round(end.Sub(accessTime).Seconds() * 1000))
			fmt.Println(accessTime, end)
			fmt.Println((end.Sub(accessTime).Seconds() * 1000))
			accessLog := model.RequestInfo{
				Timestamp: accessTime,
				Latency:   latency,
				Path:      urlPath,
				OS:        os,
			}

			fmt.Printf("%+v\n", accessLog)
			err := json.NewEncoder(w).Encode(accessLog)
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
	return http.HandlerFunc(fn)
}

// logは後に実行される必要がある。latency=0なのはおかしいから。logが最後に実行されるように色々デバッグ試してみて
