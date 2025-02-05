package main

import (
	"example/hello/clients/telegram"
	"flag"
	"log"
)

const (
	tgBotHost = "https://api.telegram.org"
)

func main() {
	tgClient := telegram.New(tgBotHost, mustToken())
	// fetcher = fetcher.New()
	// processor = processor.New()

	// consumer.Start(fetcher, processor)

}

func mustToken() string {
	token := flag.String("token-bot-token", "", "tg token access")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not exist")
	}

	return *token
}
