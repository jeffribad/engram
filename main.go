package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/engram/internal/chunk"
	"github.com/your-org/engram/internal/config"
	"github.com/your-org/engram/internal/store"
)

const version = "0.1.0"

func main() {
	// Define CLI flags
	versionFlag := flag.Bool("version", false, "Print version and exit")
	queryFlag := flag.String("query", "", "Query the engram memory store")
	ingestFlag := flag.String("ingest", "", "Ingest a file or directory into the memory store")
	configPath := flag.String("config", ".engram/config.json", "Path to config file")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("engram v%s\n", version)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize chunk store
	s, err := store.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initializing store: %v\n", err)
		os.Exit(1)
	}
	defer s.Close()

	switch {
	case *ingestFlag != "":
		if err := ingest(s, *ingestFlag); err != nil {
			fmt.Fprintf(os.Stderr, "ingest error: %v\n", err)
			os.Exit(1)
		}
	case *queryFlag != "":
		results, err := query(s, *queryFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "query error: %v\n", err)
			os.Exit(1)
		}
		for _, r := range results {
			fmt.Println(r)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}

// ingest reads content from the given path and stores it as chunks.
func ingest(s *store.Store, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %q: %w", path, err)
	}

	if info.IsDir() {
		return chunk.IngestDirectory(s, path)
	}
	return chunk.IngestFile(s, path)
}

// query searches the store for chunks relevant to the given query string.
func query(s *store.Store, q string) ([]string, error) {
	chunks, err := s.Search(q)
	if err != nil {
		return nil, err
	}

	out := make([]string, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, c.Content)
	}
	return out, nil
}
