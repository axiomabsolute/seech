package main

import (
	"bufio"
	"log"
	"os"

	"github.com/axiomabsolute/seech/index/sqlite"
	// "github.com/axiomabsolute/seech/text"
	"github.com/urfave/cli/v2"
)

// Usage:
// xsv select 2 pokemon.csv | cat -n | xargs -I {} seech index my_pokedex pokemon.csv "{}"
// seech search my_index "Bulbasaur"
// seech clear my_pokedex pokemon.csv

func addInternal(indexName string, filePath string, numberedLines []string) {
	sqlite.TrigramAddToIndex(indexName, filePath, numberedLines)
}

func add(c *cli.Context) error {
	indexName := c.Args().Get(0)
	filePath := c.Args().Get(1)
	numberedLine := c.Args().Get(2)

	log.Printf("index(%s, %s, %s)\n", indexName, filePath, numberedLine)

	sqlite.CheckAndCreate(indexName)
	addInternal(indexName, filePath, []string{numberedLine})

	return nil
}

func batch(c *cli.Context) error {
	indexName := c.Args().Get(0)
	filePath := c.Args().Get(1)
	docPath := c.Args().Get(2)

	log.Printf("batch(%s, %s, %s)\n", indexName, filePath, docPath)

	sqlite.CheckAndCreate(indexName)

	file, err := os.Open(docPath)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	batch := []string{}
	for scanner.Scan() {
		batch = append(batch, scanner.Text())
		if len(batch) == 5000 {
			addInternal(indexName, filePath, batch)
			batch = []string{}
		}
	}
	if len(batch) > 0 {
		addInternal(indexName, filePath, batch)
	}

	return nil
}

func search(c *cli.Context) error {
	indexName := c.Args().Get(0)
	query := c.Args().Get(1)

	log.Printf("query(%s, %s)\n", indexName, query)

	sqlite.TrigramSearch(indexName, query)

	return nil
}

func remove(c *cli.Context) error {
	indexName := c.Args().Get(0)
	filePath := c.Args().Get(1)

	log.Printf("remove(%s, %s)\n", indexName, filePath)

	// sqlite.RemoveFromIndex(indexName, filePath)

	return nil
}

func clear(c *cli.Context) error {
	indexName := c.Args().Get(0)

	log.Printf("clear(%s)\n", indexName)

	// sqlite.Clear(indexName)

	return nil
}

func main() {
	app := &cli.App{
		Name:  "seech",
		Usage: "Create and query full-text search indexes",
		Commands: []*cli.Command{
			{
				Name:    "trigram",
				Aliases: []string{"t"},
				Usage:   "Trigram-based indexing",
				Subcommands: []*cli.Command{
					{
						Name:      "index",
						Aliases:   []string{"i"},
						Usage:     "add index entries for file",
						UsageText: "Add entries to the specified index for the provided numbered line, which point to that line number in the given filepath. Existing entries are not removed. Lines should be numbered, e.g. with `cat -n`.",
						ArgsUsage: "index_name file_path \"numbered_line\"",
						Action:    add,
					},
					{
						Name:      "batch",
						Aliases:   []string{"b"},
						Usage:     "add index entries for file from a batch",
						UsageText: "Add entries to the specified index for the provided numbered line, which point to that line number in the given filepath. Existing entries are not removed. Lines should be numbered, e.g. with `cat -n`.",
						ArgsUsage: "index_name file_path \"numbered_line\"",
						Action:    batch,
					},
					{
						Name:      "search",
						Aliases:   []string{"s"},
						Usage:     "search index for query string",
						ArgsUsage: "index_name query",
						Action:    search,
					},
					{
						Name:      "remove",
						Aliases:   []string{"r"},
						Usage:     "remove index entries for file",
						ArgsUsage: "index_name file_path",
						Action:    remove,
					},
					{
						Name:      "clear",
						Aliases:   []string{"c"},
						Usage:     "clear all files from index",
						ArgsUsage: "index_name",
						Action:    clear,
					},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
