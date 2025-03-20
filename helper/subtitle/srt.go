package subtitle

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/asticode/go-astisub"
)

// astisub.Subtitles 这个结构看着愚笨, 实则厚实. 人生智慧也是一样, 兔子看着很快, 结果乌龟走完了全程

// DownloadSRT downloads an SRT file from a URL and parses it directly from the HTTP response body.
func DownloadSRT(ctx context.Context, url string) (subtitles *astisub.Subtitles, err error) {
	// Send HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to download SRT file", slog.String("url", url), slog.Any("error", err))
		return
	}
	defer resp.Body.Close()

	// Check for HTTP success response
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to download SRT: received status code %d", resp.StatusCode)
		slog.ErrorContext(ctx, "Non-200 response when downloading SRT", slog.String("url", url), slog.Int("status_code", resp.StatusCode))
		return
	}

	// Directly parse the SRT from HTTP response body
	subtitles, err = astisub.ReadFromSRT(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to parse SRT from response body", slog.String("url", url), slog.Any("error", err))
		return
	}

	return
}

// LoadSRT reads an SRT file and returns subtitle data.
func LoadSRT(ctx context.Context, filePath string) (subtitles *astisub.Subtitles, err error) {
	// Open the SRT file
	file, err := os.Open(filePath)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to open SRT file", slog.String("file", filePath), slog.Any("error", err))
		return
	}
	defer file.Close()

	// Parse the SRT file
	subtitles, err = astisub.ReadFromSRT(file)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to parse SRT file", slog.String("file", filePath), slog.Any("error", err))
		return
	}

	return
}

// SaveSRT writes subtitles to an SRT file.
func SaveSRT(ctx context.Context, filePath string, subtitles *astisub.Subtitles) error {
	// Open the output file
	file, err := os.Create(filePath)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create SRT file", slog.String("file", filePath), slog.Any("error", err))
		return err
	}
	defer file.Close()

	// Write subtitles to SRT format
	err = subtitles.WriteToSRT(file)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to write SRT file", slog.String("file", filePath), slog.Any("error", err))
		return err
	}

	return nil
}
