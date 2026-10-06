package cli

import (
	"testing"
)

func TestParse_defaults(t *testing.T) {
	cfg, err := Parse(nil)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Count != 1 || cfg.Delimiter != " " || cfg.ShowHelp || cfg.ShowVersion {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestParse_countOnly(t *testing.T) {
	cfg, err := Parse([]string{"7"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Count != 7 {
		t.Fatalf("count: %d", cfg.Count)
	}
}

func TestParse_flagsAfterCount(t *testing.T) {
	cfg, err := Parse([]string{"5", "-d", "-"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Count != 5 || cfg.Delimiter != "-" {
		t.Fatalf("config: %+v", cfg)
	}
}

func TestParse_flagsBeforeCount(t *testing.T) {
	cfg, err := Parse([]string{"-d", "_", "3"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Count != 3 || cfg.Delimiter != "_" {
		t.Fatalf("config: %+v", cfg)
	}
}

func TestParse_delimiterEquals(t *testing.T) {
	cfg, err := Parse([]string{"--delimiter=::", "2"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Delimiter != "::" || cfg.Count != 2 {
		t.Fatalf("config: %+v", cfg)
	}
}

func TestParse_delimiterShortAttached(t *testing.T) {
	cfg, err := Parse([]string{"-d_", "1"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Delimiter != "_" {
		t.Fatalf("delimiter: %q", cfg.Delimiter)
	}
}

func TestParse_emptyDelimiter(t *testing.T) {
	cfg, err := Parse([]string{"-d", "", "2"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Delimiter != "" || cfg.Count != 2 {
		t.Fatalf("config: %+v", cfg)
	}
}

func TestParse_helpAndVersion(t *testing.T) {
	cfg, err := Parse([]string{"--help", "--version"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.ShowHelp || !cfg.ShowVersion {
		t.Fatalf("config: %+v", cfg)
	}
}

func TestParse_shortHelpAndVersion(t *testing.T) {
	cfg, err := Parse([]string{"-h"})
	if err != nil || !cfg.ShowHelp {
		t.Fatalf("config: %+v err=%v", cfg, err)
	}
	cfg, err = Parse([]string{"-v"})
	if err != nil || !cfg.ShowVersion {
		t.Fatalf("config: %+v err=%v", cfg, err)
	}
}

func TestParse_unknownFlag(t *testing.T) {
	_, err := Parse([]string{"--nope"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParse_delimiterMissingValue(t *testing.T) {
	_, err := Parse([]string{"-d"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParse_tooManyPositionals(t *testing.T) {
	_, err := Parse([]string{"1", "2"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParse_invalidCount(t *testing.T) {
	_, err := Parse([]string{"abc"})
	if err == nil {
		t.Fatal("expected error")
	}
	_, err = Parse([]string{"0"})
	if err == nil {
		t.Fatal("expected error")
	}
}
