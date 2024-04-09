package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(todoDB)

	server := &http.Server{ // httpServerの作成
		Addr:    port, // デフォは:8080
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill) // 割り込みと強制終了のシグナル
	defer stop()

	wg := &sync.WaitGroup{} // goroutineの終了を待つ
	wg.Add(1)

	go func() {
		defer wg.Done() // goroutineが完了したらwgを1減らす(deferの処理は下から上)
		<-ctx.Done()    // ctxがキャンセルされるまで以下の処理は実行されないようにする(中断処理信号が受信"されたら"サーバーをシャットダウン)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			log.Println("graceful shutdown miss", err)

		} else {
			log.Println("graceful shutdown")
		}
	}()

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	wg.Wait()
	return nil
}
