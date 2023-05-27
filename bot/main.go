package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"greenfield-deploy/bot/config"

	"golang.org/x/sync/semaphore"
)

var MaxWorkersNum = int64(10)

func main() {
	conf := config.Load("config/config.yaml")

	bot, err := NewBot(conf)
	if err != nil {
		log.Fatal(err)
	}

	updates, err := bot.Channel()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("authorized on account %s", bot.UserName())

	ctx, cancel := context.WithCancel(context.Background())

	go initShutdown(cancel)

	// workerpool
	sm := semaphore.NewWeighted(MaxWorkersNum)
	var wg sync.WaitGroup
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			if err := sm.Acquire(ctx, 1); err != nil {
				log.Println(err)
				continue
			}
			wg.Add(1)
			go worker(sm, &wg, bot, update)
		}
	}
}

func initShutdown(cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGTERM)
	log.Println("shutdown by signal", <-signalChan)
	cancel()
}
