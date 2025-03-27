package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Nls0n/mattermost-poll-bot/internal/config"
	"github.com/Nls0n/mattermost-poll-bot/internal/handlers"
	"github.com/Nls0n/mattermost-poll-bot/internal/logger"
	"github.com/Nls0n/mattermost-poll-bot/internal/storage"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/sirupsen/logrus"
)

func main() {
	// Инициализация логгера
	logger.Init()
	logger.Log.Info("Starting Mattermost Poll Bot")

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to load configuration")
	}

	// Инициализация Tarantool
	tarantoolStorage, err := storage.New(cfg.TarantoolAddr)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err,
			"addr":  cfg.TarantoolAddr,
		}).Fatal("Failed to connect to Tarantool")
	}
	defer tarantoolStorage.Close()

	// Клиент Mattermost
	mmClient := model.NewAPIv4Client(cfg.MattermostURL)
	mmClient.SetToken(cfg.MattermostToken)

	// Проверка подключения
	if _, resp := mmClient.GetMe(""); resp.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": resp.Error,
			"url":   cfg.MattermostURL,
		}).Fatal("Failed to connect to Mattermost")
	}

	// WebSocket клиент
	wsClient, err := model.NewWebSocketClient4(cfg.MattermostURL, cfg.MattermostToken)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to create WebSocket client")
	}
	defer wsClient.Close()

	wsClient.Listen()
	logger.Log.Info("Bot is now listening for events")

	// Обработка событий
	go func() {
		for event := range wsClient.EventChannel {
			if event.EventType() == model.WebsocketEventPosted {
				handlers.ProcessEvent(mmClient, tarantoolStorage, event)
			}
		}
	}()

	// Graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	logger.Log.Info("Received shutdown signal")
}
