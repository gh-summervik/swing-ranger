package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/summervik/swing-ranger/internal/config"
	"github.com/summervik/swing-ranger/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	cfg, err := parseArgs(os.Args)
	if err != nil {
		fmt.Println(err)
		fmt.Println()
		showUsage()
		os.Exit(1)
	}

	if cfg.ShowHelp {
		showUsage()
		os.Exit(0)
	}

	if strings.TrimSpace(cfg.Command) == "" {
		fmt.Println("A command is required.")
		showUsage()
		os.Exit(2)
	}

	if strings.TrimSpace(cfg.Command) == "init" {
		msgs := initApp()

		fmt.Println()

		for _, msg := range msgs {
			fmt.Println(msg)
		}

		fmt.Println()
		fmt.Println(strings.Repeat("-", 20))
		fmt.Println()
		showUsage()
		os.Exit(0)
	}

	secrets, err := config.LoadSecrets()
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	cfg.Secrets = secrets

	appconfig, err := config.LoadAppConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}
	cfg.AppConfig = appconfig

	commsSvc := service.NewCommsService(cfg)
	dbSvc, err := service.NewDbService(cfg, commsSvc)
	if err != nil {
		commsSvc.Communicate(fmt.Sprintf("Failed to create database service: %v", err))
		os.Exit(4)
	}

	defer dbSvc.Command.Close()
	defer dbSvc.Query.Close()

	switch cfg.Command {
	case "test":
		if err := doTest(cfg, dbSvc); err != nil {
			fmt.Println(err)
		}
	case "collect":
		collectSymbols(cfg, commsSvc, dbSvc)
	case "backtest":
		if len(cfg.Data) == 0 {
			commsSvc.Communicate("Missing name of backtest.")
			os.Exit(1)
		}
		if err := doBacktest(cfg, dbSvc); err != nil {
			fmt.Println(err)
		}
	}

	os.Exit(0)
}