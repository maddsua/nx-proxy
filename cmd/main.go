package main

import (
	"context"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	nxproxy "github.com/maddsua/nx-proxy"
	"github.com/maddsua/nx-proxy/api_models"
)

func main() {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	configFileEntries := LoadConfigFile()
	if configFileEntries == nil {
		slog.Warn("No config file found")
	}

	var client nxproxy.Client

	if val, ok := GetConfigOpt(configFileEntries, "SECRET_TOKEN"); ok {
		token, err := nxproxy.ParseServerToken(val)
		if err != nil {
			slog.Error("STARTUP: Parse secret token",
				slog.String("err", err.Error()))
			os.Exit(1)
		}
		client.Token = token
	} else {
		slog.Warn("STARTUP: Secret token not provided")
	}

	if val, ok := GetConfigOpt(configFileEntries, "AUTH_URL"); ok {

		url, err := url.Parse(val)
		if err != nil {
			slog.Error("STARTUP: Parse auth server url",
				slog.String("err", err.Error()))
			os.Exit(1)
		}
		client.URL = url

	} else {
		slog.Error("STARTUP: Auth server url not provided")
		os.Exit(1)
	}

	runID := uuid.New()
	runAt := time.Now()

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := NewServiceHub()

	metricsTicker := time.NewTicker(30 * time.Second)
	tableTicker := time.NewTicker(15 * time.Second)

	wg.Add(1)

	go func() {

		defer wg.Done()

		var doUpdate = func() {

			//	todo: pull data from the controller

			metrics := api_models.Metrics{
				Deltas: make([]api_models.Delta, 0),
				Service: api_models.Service{
					RunID:  runID,
					Uptime: int64(time.Since(runAt).Seconds()),
				},
			}

			if err := client.PostMetrics(&metrics); err != nil {
				slog.Error("API: PostMetrics",
					slog.String("err", err.Error()))
				return
			}

			slog.Debug("API: PostMetrics OK")
		}

		doneCh := ctx.Done()

		for {
			select {
			case <-metricsTicker.C:
				doUpdate()
			case <-doneCh:
				doUpdate()
				return
			}
		}
	}()

	go func() {

		var pullTable = func() {

			table, err := client.PullTable()
			if err != nil {
				slog.Error("API: PullTable",
					slog.String("err", err.Error()))
				return
			}

			if err := hub.ApplySlots(table.Slots); err != nil {
				slog.Error("API: PullTable",
					slog.String("err", err.Error()))
				return
			}

			slog.Debug("API: PullTable OK")
		}

		doneCh := ctx.Done()

		for {

			pullTable()

			select {
			case <-tableTicker.C:
				continue
			case <-doneCh:
				return
			}
		}
	}()

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)

	<-exitCh
	cancel()
	hub.Close()

	wg.Wait()
}
