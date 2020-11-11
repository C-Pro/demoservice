package main

import (
	"context"
	"demoservice/pkg/service"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	log.Println("service is started")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(ctx context.Context) {
		defer wg.Done()
		if err := service.Run(ctx); err != nil {
			log.Fatalf("service start failed: %v", err)
		}
	}(ctx)

	<-ch
	log.Println("service is shutting down...")
	cancel()
	wg.Wait()
	log.Println("shut down finished")
}
