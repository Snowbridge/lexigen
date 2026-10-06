package main

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"

	"github.com/Snowbridge/lexigen/internal/db"
)

//go:embed words.sqlite
var wordsSQLite []byte

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [count]\n", os.Args[0])
}

func main() {
	count := 1
	switch len(os.Args) - 1 {
	case 0:
	case 1:
		n, err := strconv.Atoi(os.Args[1])
		if err != nil || n < 1 {
			fmt.Fprintln(os.Stderr, "count must be a positive integer")
			usage()
			os.Exit(1)
		}
		count = n
	default:
		usage()
		os.Exit(1)
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
	if count > max {
		fmt.Fprintf(os.Stderr, "requested %d phrases, maximum is %d\n", count, max)
		os.Exit(1)
	}

	phrases, err := store.Phrases(count)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, phrase := range phrases {
		fmt.Println(phrase)
	}
}
