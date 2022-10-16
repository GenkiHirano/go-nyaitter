package server

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/GenkiHirano/go-nyaitter.git/twitter"
	"github.com/labstack/echo"
)

func RunAPIServer() {
	e := echo.New()
	e.GET("/auth", twitter.AuthTwitter)
	e.GET("/callback", twitter.Callback)

	// サーバー開始
	go func() {
		if err := e.Start(":3022"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
