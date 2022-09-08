package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tinderutf/api"
	"tinderutf/db"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db.StartDB()

	server := api.NewServer(ctx, ":51000", false)

	srv := &http.Server{
		Addr:    server.Port(),
		Handler: server.Start(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	db.Close()

	ctx, cancel = context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
