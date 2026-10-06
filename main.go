package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/Snowbridge/lexigen/internal/cli"
	"github.com/Snowbridge/lexigen/internal/db"
)

//go:embed words.sqlite
var wordsSQLite []byte

func main() {
	prog := os.Args[0]
	cfg, err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprint(os.Stderr, cli.Help(prog))
		os.Exit(1)
	}

	if cfg.ShowHelp {
		fmt.Print(cli.Help(prog))
		return
	}
	if cfg.ShowVersion {
		fmt.Println(cli.VersionLine(Version))
		return
	}

	store, err := db.Open(wordsSQLite)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer store.Close()

	max, err := store.MaxPhrases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if cfg.Count > max {
		fmt.Fprintf(os.Stderr, "requested %d phrases, maximum is %d\n", cfg.Count, max)
		os.Exit(1)
	}

	phrases, err := store.Phrases(cfg.Count, cfg.Delimiter)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, phrase := range phrases {
		fmt.Println(phrase)
	}
}
