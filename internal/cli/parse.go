package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const DefaultDelimiter = "-"

type Config struct {
	Count       int
	Delimiter   string
	ShowHelp    bool
	ShowVersion bool
}

func Parse(args []string) (Config, error) {
	cfg := Config{
		Count:     1,
		Delimiter: DefaultDelimiter,
	}

	var positionals []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--help", a == "-h":
			cfg.ShowHelp = true
		case a == "--version", a == "-v":
			cfg.ShowVersion = true
		case a == "-d", a == "--delimiter":
			if i+1 >= len(args) {
				return Config{}, errors.New("delimiter flag requires a value")
			}
			i++
			cfg.Delimiter = args[i]
		case strings.HasPrefix(a, "--delimiter="):
			cfg.Delimiter = strings.TrimPrefix(a, "--delimiter=")
		case strings.HasPrefix(a, "-d") && len(a) > 2:
			cfg.Delimiter = a[2:]
		case strings.HasPrefix(a, "-"):
			return Config{}, fmt.Errorf("unknown flag: %s", a)
		default:
			positionals = append(positionals, a)
		}
	}

	if len(positionals) > 1 {
		return Config{}, errors.New("too many arguments")
	}
	if len(positionals) == 1 {
		n, err := strconv.Atoi(positionals[0])
		if err != nil || n < 1 {
			return Config{}, errors.New("count must be a positive integer")
		}
		cfg.Count = n
	}

	return cfg, nil
}

func Help(prog string) string {
	return fmt.Sprintf(`Lexigen
Print random English adjective–noun lexemes.

Usage:
  %s [flags] [count]

Flags:
  -d, --delimiter <string>  Text between adjective and noun (default: "%s")
  -h, --help                Show this help
  -v, --version             Show version

count is the number of lexemes to print (default: 1).
`, prog, DefaultDelimiter)
}

func VersionLine(version string) string {
	return "lexigen " + version
}
