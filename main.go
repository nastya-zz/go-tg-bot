package main

import (
	tgClient "example/hello/clients/telegram"
	event_consumer "example/hello/consumer/event-consumer"
	"example/hello/events/telegram"
	"example/hello/storage/files"
	"flag"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.New(tgClient.New(tgBotHost, mustToken()), files.New(storagePath))
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
