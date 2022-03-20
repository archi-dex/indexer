package main

import (
	"sync"

	"github.com/hans-m-song/archidex/indexer/pkg/index"
	"github.com/hans-m-song/archidex/indexer/pkg/parse"
	"github.com/hans-m-song/archidex/indexer/pkg/submit"
	"github.com/hans-m-song/archidex/indexer/pkg/util"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "indexer",
	Short: "recursively index a directory and batch submit",
	Run:   run,
}

func init() {
	util.InitFlags(cmd)
}

func main() {
	cmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	options := util.InitOptions(cmd)
	logger := util.InitLogger(options.Debug)
	defer logger.Sync()

	files := make(chan string)
	entities := make(chan parse.Entity)
	wg := sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(files)

		if err := index.IndexFromWalk(logger, files, options.Cwd, options.SearchDir, options.Ignore, options.Match); err != nil {
			logger.Fatal(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(entities)

		if err := parse.Parser(logger, files, entities, options.Pattern); err != nil {
			logger.Fatal(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := submit.Submitter(logger, entities, options.OutputFile, options.OutputEndpoint, options.DryRun); err != nil {
			logger.Fatal(err)
		}
	}()
}
