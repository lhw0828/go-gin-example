package main

import (
	"context"
	"fmt"
	"github.com/lhw0828/go-gin-example/pkg/logging"
	"github.com/lhw0828/go-gin-example/pkg/setting"
	"github.com/lhw0828/go-gin-example/routers"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			logging.Fatal("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logging.Info("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		logging.Fatal("Server Shutdown:", err)
	}

	logging.Info("Server exiting")

	//err := s.ListenAndServe()
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	//endless.DefaultReadTimeOut = setting.ReadTimeout
	//endless.DefaultWriteTimeOut = setting.WriteTimeout
	//endless.DefaultMaxHeaderBytes = 1 << 20
	//endPoint := fmt.Sprintf(":%d", setting.HttpPort)
	//
	//server := endless.NewServer(endPoint, routers.InitRouter())
	//server.BeforeBegin = func(add string) {
	//	log.Printf("Actual pid is %d", syscall.Getpid())
	//}
	//
	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Printf("Server err: %v", err)
	//}
}
