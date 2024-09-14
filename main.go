package main

import (
	"fmt"
	"github.com/lhw0828/go-gin-example/pkg/setting"
	"github.com/lhw0828/go-gin-example/routers"
	"log"
	"net/http"
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

	err := s.ListenAndServe()
	if err != nil {
		log.Println(err)
		return
	}
}
