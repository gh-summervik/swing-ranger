package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/summervik/swing-ranger/internal/config"
	"github.com/summervik/swing-ranger/internal/service"

	_ "github.com/lib/pq"
)

type CommandArg struct {
	Command     string
	Description string
}

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

	// If the user asked for help, just provide the help and exit.
	if cfg.ShowHelp {
		showUsage()
		os.Exit(0)
	}

	// The absense of a command ends the process.
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

	// Extract the secrets from the super secret json file.
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

	// establish services.
	commsSvc := service.NewCommsService(cfg)
	dbSvc, err := service.NewDbService(cfg, commsSvc)
	if err != nil {
		commsSvc.Communicate(fmt.Sprintf("Failed to create database service: %v", err))
		os.Exit(4)
	}

	defer dbSvc.Command.Close()
	defer dbSvc.Query.Close()

	// process the specified command.
	switch cfg.Command {
	case "test":
		doTest(cfg)
	case "collect":
		collectSymbols(cfg, commsSvc, dbSvc)
	}

	os.Exit(0)
}

func doTest(cfg config.Config) {
	data, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(data))
}
