package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	searchpb "yadro.com/course/proto/search"
	"yadro.com/course/search/adapters/db"
	"yadro.com/course/search/adapters/eventbus"
	searchgrpc "yadro.com/course/search/adapters/grpc"
	"yadro.com/course/search/adapters/initiator"
	"yadro.com/course/search/adapters/words"
	"yadro.com/course/search/config"
	"yadro.com/course/search/core"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := mustMakeLogger(cfg.LogLevel)

	if err := run(cfg, log); err != nil {
		log.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func run(cfg config.Config, log *slog.Logger) error {
	log.Info("starting search server")

	repo, err := db.New(log, cfg.DBAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	wordsClient, err := words.NewClient(cfg.WordsAddress, log)
	if err != nil {
		return fmt.Errorf("failed to create words client: %w", err)
	}

	svc := core.NewService(log, repo, wordsClient)

	initiatorAdapter := initiator.NewInitiator(log, svc, cfg.IndexTTL)

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", cfg.Address, err)
	}

	subscriber, err := eventbus.NewSubscriber(cfg.BrokerAddress, log, svc)
	if err != nil {
		return fmt.Errorf("failed to create eventbus subscriber: %w", err)
	}
	defer subscriber.Close()

	if err := subscriber.Subscribe(context.Background()); err != nil {
		return fmt.Errorf("failed to subscribe to events: %w", err)
	}

	s := grpc.NewServer()
	searchpb.RegisterSearchServer(s, searchgrpc.NewServer(svc))
	reflection.Register(s)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	initiatorAdapter.Start(ctx)

	go func() {
		<-ctx.Done()
		log.Info("shutting down server")
		s.GracefulStop()
	}()

	log.Info("server listening", "address", cfg.Address)
	if err := s.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func mustMakeLogger(logLevel string) *slog.Logger {
	var level slog.Level
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
