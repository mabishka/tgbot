package main

import (
	"context"
	"fmt"
	"log"
	"tgbot/internal/config"
	"tgbot/internal/logger"
	"tgbot/internal/repository/chat"
	"tgbot/internal/repository/db"
	"tgbot/internal/repository/speach"
	"tgbot/internal/service/bot"
	"tgbot/internal/service/chatservice"
	"tgbot/internal/service/speachservice"

	"golang.org/x/sync/errgroup"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

const chatStreamCount = 1
const speachStreamCount = 1
const queueSpeachSize = 100

func main() {
	fmt.Println("Build version: ", buildVersion)
	fmt.Println("Build date: ", buildDate)
	fmt.Println("Build commit: ", buildCommit)

	if err := create(context.WithCancelCause(context.Background())); err != nil {
		log.Fatalf("exist with error: %v", err)
	}
}

func create(ctx context.Context, fnCancel context.CancelCauseFunc) error {

	config := config.New()
	if err := logger.InitLogger(config.GetLogLevel()); err != nil {
		panic(err)
	}

	dbConn, err := db.New(config.GetConnectionString())
	if err != nil {
		return err
	}
	speachConn := speach.New(config.GetSpeachAuthHost(), config.GetSpeachRQUID(), config.GetSpeachAuthKey())
	speachService := speachservice.New(
		speachservice.WithSpeach(speachConn),
		speachservice.WithHost(config.GetSpeachRequestHost()),
	)

	chatConn := chat.New(config.GetChatAuthHost(), config.GetChatRQUID(), config.GetChatAuthKey())
	chatService := chatservice.New(config.GetChatRequestHost(), chatConn)

	botService := bot.New(config.GetBotToken(), dbConn, speachService, chatService)

	wg, sendCtx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		return speachService.Start(sendCtx, speachStreamCount, queueSpeachSize)
	})
	wg.Go(func() error {
		return botService.Start(sendCtx, speachStreamCount, queueSpeachSize)
	})

	return wg.Wait()
}
