package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/summervik/swing-ranger/internal/config"

	_ "github.com/lib/pq"
)

type CommandArg struct {
	Command     string
	Description string
}

func parseArgs(args []string) (config.Config, error) {
	cfg := config.Config{
		Command:  "",
		Verbose:  false,
		ShowHelp: false,
		Data:     nil,
		Secrets:  nil,
	}

	count := len(args)

	for i := 1; i < count; i++ {
		arg := strings.ToLower(os.Args[i])
		switch arg {
		case "init", "update", "test":
			cfg.Command = arg
		case "backtest", "watchlist":
			cfg.Command = arg
			i++
			if i >= count {
				return cfg, fmt.Errorf("Expecting a backtest name after %s", args[i-1])
			}
			cfg.Data = make([]string, 0, 1)
			cfg.Data = append(cfg.Data, args[i])
		case "collect":
			cfg.Command = arg
			i++
			if i >= count {
				return cfg, fmt.Errorf("Expecting one or more symbols after %s", args[i-1])
			}

			symbolsStr := strings.TrimSpace(args[i])
			if symbolsStr == "" {
				return cfg, fmt.Errorf("Expecting symbol(s) after %s\n", args[i-1])
			}
			rawSymbols := strings.Split(symbolsStr, ",")
			cfg.Data = make([]string, 0, len(rawSymbols))
			for _, s := range rawSymbols {
				trimmed := strings.TrimSpace(s)
				if trimmed != "" {
					cfg.Data = append(cfg.Data, strings.ToUpper(trimmed))
				}
			}
			if len(cfg.Data) == 0 {
				return cfg, fmt.Errorf("No valid symbols provided")
			}

		case "-v", "--verbose":
			cfg.Verbose = true
		case "-h", "help", "--help", "-?", "?":
			cfg.ShowHelp = true
		default:
			return cfg, fmt.Errorf("Unknown argument: %s\n", os.Args[i])
		}
	}

	return cfg, nil
}

func showUsage() {
	cmds := []CommandArg{
		{"[init]", "Command to check basics before first-time use. Will NOT injure existing initializations."},
		{"[collect <S1,S2,S3>]", "Command to capture price history for provided symbols. Use commas to separate symbols; no spaces."},
		{"[update]", "Command to update price action data for existing symbols (symbols previously captured)."},
		{"[backtest <name>]", "Run a configured backtest (e.g. squeeze_breakout)"},
		{"[watchlist <path>]", "Reads in a CSV file for setting up watchlists (e.g., SPY)"},
		{"[-h | -? | ? | --help]", "Show help"},
		{"[-v | --verbose]", "Verbose output"},
	}

	allArgs := make([]string, len(cmds))
	for i, cmd := range cmds {
		allArgs[i] = cmd.Command
	}

	var colSize = 0
	for i := 0; i < len(cmds); i++ {
		if len(cmds[i].Command) > colSize {
			colSize = len(cmds[i].Command)
		}
	}

	var exeName = filepath.Base(os.Args[0])
	fmt.Printf("\n%s\t\tSwing Ranger - a trader's toolbox\n", exeName)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Printf("%s %s\n", exeName, strings.Join(allArgs, " "))
	fmt.Println()

	for i := 0; i < len(cmds); i++ {
		fmt.Printf("%s\n", fmt.Sprintf("%-*s\t%s", colSize, cmds[i].Command, cmds[i].Description))
	}
	fmt.Println()
	fmt.Println("Only one command (e.g., init, collect, update, backtest, watchlist, etc.) is allowed; last one wins.")
	fmt.Println()
	fmt.Println("Examples")
	fmt.Println(strings.Repeat("-", 20))
	fmt.Println()
	fmt.Printf("`%s collect MSFT,TSLA`", exeName)
	fmt.Println("\tWill collect and preserve historical price information for MSFT and TSLA.")
	fmt.Println()
	fmt.Printf("`%s collect MSFT,TSLA -v`", exeName)
	fmt.Println("\tCollects historical price information for MSFT and TSLA with verbose output.")
	fmt.Println()
	fmt.Printf("`%s update`", exeName)
	fmt.Println("\tUpdates price action for existing symbols.")
	fmt.Println()
	fmt.Printf("`%s backtest squeeze_breakout`", exeName)
	fmt.Println("\tRuns the squeeze breakout backtest on SPY.")
	fmt.Println()
	fmt.Printf("`%s watchlist ./data/watchlist.csv`", exeName)
	fmt.Println("\tReplaces the watchlists in the database with those provided by the CSV.")
	fmt.Println()
	fmt.Println(strings.Repeat("-", 20))
	fmt.Println()
	fmt.Println("A secrets.json file is required; here is an example of the content.")
	fmt.Println()
	fmt.Println(`{
  "ConnectionStrings": {
    "Command": "Your connection string here",
    "Query": "Your connection string here"
  }
}`)
	fmt.Println()
	fmt.Println("A config.json can be used to configure chart and backtesting parameters. Here is an example.")
	fmt.Println()
	fmt.Println(`{
    "Chart": {
        "MovingAverages": [
            "21SC",
            "55SC",
            "233SC"
        ],
        "BollingerBandsMovingAverage": "20SC",
        "MACD": "12,26,9C",
        "RSI": "14C",
        "Backtests": {
            "squeeze_breakout": {
                "type": "squeeze_breakout",
                "squeezeLookback": 50,
                "minSqueezeBars": 3,
                "minRSI": 50
            }
        }
    }
}`)
}
