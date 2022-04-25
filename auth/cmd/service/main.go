package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ypxd/diplom/auth/internal/repository"
	"github.com/Ypxd/diplom/auth/internal/repository/redis"
	"github.com/Ypxd/diplom/auth/internal/server"
	"github.com/Ypxd/diplom/auth/internal/service"
	"github.com/Ypxd/diplom/auth/internal/transport/http"
	"github.com/Ypxd/diplom/auth/utils"
	"github.com/Ypxd/diplom/shared"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	shared.InitLogger(utils.GetConfig().Logger)

	redisCon, err := redis.Connect()
	if err != nil {
		err = errors.New(fmt.Sprintf("error occurre while connection to cache: %s", err.Error()))
		shared.GetLogger().Fatal(err)
	}

	repo, conn, err := repository.NewRepo(redisCon)
	if err != nil {
		shared.GetLogger().Fatal(err)
	}

	services := service.NewService(repo, conn)

	handlers := handler.InitHandlers(services)
	srv := server.NewServer(handlers)

	stopCh := make(chan int)
	go startHttpApi(stopCh, srv)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	stopCh <- 1

	<-stopCh
	shared.GetLogger().Log(4, "Graceful Shutdown Service")

}

func startHttpApi(stopChan chan int, srv *server.Server) {
	conf := utils.GetConfig()
	go func() {
		shared.GetLogger().Infof("HTTP API server listeing at: %s:%d",
			conf.Server.Host,
			conf.Server.Port)

		err := srv.Run()
		if !errors.Is(err, http.ErrServerClosed) {
			shared.GetLogger().Fatalf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	<-stopChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		shared.GetLogger().Fatalf("Server forced to shutdown: %s", err.Error())
	}
	shared.GetLogger().Log(4, "HTTP API service stopped")
	stopChan <- 1
}
