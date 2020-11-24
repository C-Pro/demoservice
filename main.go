package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"demoservice/pkg/db"
	"demoservice/pkg/service"
)

func main() {
	log.Println("service is started")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(ctx context.Context) {
		defer wg.Done()
		if err := service.Run(ctx, conn); err != nil {
			log.Fatalf("service start failed: %v", err)
		}
	}(ctx)

	<-ch
	log.Println("service is shutting down...")
	cancel()
	wg.Wait()
	log.Println("shut down finished")
}
