package main

import (
	"github.com/lhw0828/go-gin-example/models"
	"github.com/lhw0828/go-gin-example/pkg/logging"
	"github.com/robfig/cron"
	"time"
)

func main() {
	logging.Info("Starting...")

	c := cron.New()

	errTag := c.AddFunc("* * * * * *", func() {
		logging.Info("Run models.CleanAllTag...")
		models.CleanAllTag()
	})
	if errTag != nil {
		return
	}
	errArticle := c.AddFunc("* * * * * *", func() {
		logging.Info("Run models.CleanAllArticles...")
		models.CleanAllArticle()
	})
	if errArticle != nil {
		return
	}

	c.Start()

	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}

}
