package fontsource

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func DownloadFonts(ctx context.Context, logger *slog.Logger, outputDir, formats, weights, styles, subsets string) error {
	logger.InfoContext(ctx, "starting font source downloader")

	fonts, err := fetchFonts(ctx, logger)
	if err != nil {
		logger.ErrorContext(ctx, "an error occurred", slog.String("error", err.Error()))
		return err
	}

	successfullyDownloadedFonts := []Font{}
	for _, font := range fonts {
		successfullyDownloadedFonts = append(
			successfullyDownloadedFonts,
			downloadFont(ctx, logger, outputDir, formats, weights, styles, subsets, font),
		)
	}

	logger.InfoContext(ctx, "font source downloader finished")

	if len(successfullyDownloadedFonts) > 0 {
		logger.InfoContext(ctx, "successfully downloaded fonts", slog.Int("count", len(successfullyDownloadedFonts)))

		// Store successfully downloaded fonts in the output directory
		filename := filepath.Join(outputDir, "fonts.json")
		file, err := os.Create(filename)
		if err != nil {
			logger.ErrorContext(ctx, "an error occurred", slog.String("error", err.Error()))
			return err
		}
		defer file.Close()

		if err := json.NewEncoder(file).Encode(successfullyDownloadedFonts); err != nil {
			logger.ErrorContext(ctx, "an error occurred", slog.String("error", err.Error()))
			return err
		}

		logger.InfoContext(ctx, "successfully stored successfully downloaded fonts", slog.String("filename", filename))
	}

	return nil
}

type Font struct {
	ID           string   `json:"id"`
	Family       string   `json:"family"`
	Subsets      []string `json:"subsets"`
	Weights      []int    `json:"weights"`
	Styles       []string `json:"styles"`
	DefSubset    string   `json:"defSubset"`
	Variable     bool     `json:"variable"`
	LastModified string   `json:"lastModified"`
	Category     string   `json:"category"`
	License      string   `json:"license"`
	Type         string   `json:"type"`
}

func fetchFonts(ctx context.Context, logger *slog.Logger) ([]Font, error) {
	logger.InfoContext(ctx, "fetching fonts")

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Get("https://api.fontsource.org/v1/fonts")
	if err != nil {
		logger.ErrorContext(ctx, "an error occurred", slog.String("error", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.ErrorContext(ctx, "an error occurred", slog.String("error", resp.Status))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var fonts []Font
	if err := json.NewDecoder(resp.Body).Decode(&fonts); err != nil {
		logger.ErrorContext(ctx, "an error occurred", slog.String("error", err.Error()))
		return nil, err
	}

	logger.InfoContext(ctx, "downloaded fonts", slog.Int("count", len(fonts)))
	return fonts, nil
}

func downloadFont(ctx context.Context, logger *slog.Logger, outputDir, formats, weights, styles, subsets string, font Font) Font {
	logger.InfoContext(ctx, "downloading font", slog.String("id", font.ID))

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	splitFormats := strings.Split(formats, ",")
	splitStyles := strings.Split(styles, ",")
	splitSubsets := strings.Split(subsets, ",")

	var splitWeights []int
	for _, weight := range strings.Split(weights, ",") {
		convertedWeight, err := strconv.Atoi(weight)
		if err != nil {
			logger.ErrorContext(ctx, "an error occurred", slog.String("error", err.Error()))
			continue
		}

		splitWeights = append(splitWeights, convertedWeight)
	}

	sucessfullyDownloadedFont := Font{}

	for _, format := range splitFormats {
		for _, weight := range splitWeights {
			for _, style := range splitStyles {
				for _, subset := range splitSubsets {
					url := fmt.Sprintf("https://cdn.jsdelivr.net/fontsource/fonts/%s@latest/%s-%d-%s.%s", font.ID, subset, weight, style, format)

					resp, err := httpClient.Get(url)
					if err != nil {
						logger.ErrorContext(ctx, "an error occurred", slog.String("error", err.Error()))
						continue
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						logger.ErrorContext(ctx, "an error occurred", slog.String("error", resp.Status))
						continue
					}

					filename := fmt.Sprintf("%s-%s-%d-%s.%s", font.ID, subset, weight, style, format)
					outputPath := filepath.Join(outputDir, filename)

					body, err := io.ReadAll(resp.Body)
					if err != nil {
						logger.ErrorContext(ctx, "an error occurred", slog.String("error", err.Error()))
						continue
					}

					if err := os.WriteFile(outputPath, body, 0o644); err != nil {
						logger.ErrorContext(ctx, "an error occurred", slog.String("error", err.Error()))
						continue
					}

					logger.InfoContext(ctx, "downloaded font", slog.String("filename", filename))

					if sucessfullyDownloadedFont.ID == "" {
						sucessfullyDownloadedFont = Font{
							ID:           font.ID,
							Family:       font.Family,
							Subsets:      []string{subset},
							Weights:      []int{weight},
							Styles:       []string{style},
							DefSubset:    font.DefSubset,
							Variable:     font.Variable,
							LastModified: font.LastModified,
							Category:     font.Category,
							License:      font.License,
							Type:         font.Type,
						}
					} else {
						sucessfullyDownloadedFont.Subsets = append(sucessfullyDownloadedFont.Subsets, subset)
						sucessfullyDownloadedFont.Weights = append(sucessfullyDownloadedFont.Weights, weight)
						sucessfullyDownloadedFont.Styles = append(sucessfullyDownloadedFont.Styles, style)
					}
				}
			}
		}
	}

	return sucessfullyDownloadedFont
}
