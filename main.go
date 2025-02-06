package main

import (
	"context"
	tgClient "example/hello/clients/telegram"
	event_consumer "example/hello/consumer/event-consumer"
	"example/hello/events/telegram"
	"example/hello/storage/sqlite"
	"flag"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
	//storagePath = "files_storage"
	batchSize  = 100
	sqlitePath = "data/sqlite/storage.db"
)

func main() {
	//s := files.New(storagePath)

	s, err := sqlite.New(sqlitePath)
	if err != nil {
		log.Fatal("can't init database", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init database", err)
	}

	eventsProcessor := telegram.New(tgClient.New(tgBotHost, mustToken()), s)
	log.Println("Application started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "tg token access")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not exist")
	}

	return *token
}
