package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	// This controls the maxprocs environment variable in container runtimes.
	// see https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/ZerGo0/fontsourcedownloader/internal/log"
	"github.com/ZerGo0/fontsourcedownloader/pkg/services/fontsource"
)

func main() {
	outputDir := flag.String("out", "", "output directory")
	formats := flag.String("formats", "woff2,woff", "font formats comma separated")
	weights := flag.String("weights", "400", "font weights comma separated")
	styles := flag.String("styles", "normal", "font styles comma separated")
	subsets := flag.String("subsets", "latin", "font subsets comma separated")
	flag.Parse()

	if *outputDir == "" || *formats == "" || *weights == "" || *styles == "" || *subsets == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Logger configuration
	logger := log.New(
		log.WithLevel(os.Getenv("LOG_LEVEL")),
		log.WithSource(),
	)

	if err := run(logger, *outputDir, *formats, *weights, *styles, *subsets); err != nil {
		logger.ErrorContext(context.Background(), "an error occurred", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run(logger *slog.Logger, outputDir, formats, weights, styles, subsets string) error {
	ctx := context.Background()

	_, err := maxprocs.Set(maxprocs.Logger(func(s string, i ...interface{}) {
		logger.DebugContext(ctx, fmt.Sprintf(s, i...))
	}))
	if err != nil {
		return fmt.Errorf("setting max procs: %w", err)
	}

	if err := fontsource.DownloadFonts(ctx, logger, outputDir, formats, weights, styles, subsets); err != nil {
		return err
	}

	return nil
}
