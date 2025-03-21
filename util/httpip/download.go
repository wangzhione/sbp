package httpip

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/wangzhione/sbp/chain"
)

// Download 下载 uri 到本地文件 outputPath
func Download(ctx context.Context, uri, outputpath string, headerargs ...map[string]string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		slog.ErrorContext(ctx, "http.NewRequestWithContext error", "error", err, "uri", uri)
		return err
	}

	// 设置默认 req header
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set(chain.XRquestID, chain.GetTraceID(ctx))
	for _, headers := range headerargs {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := HTTPClient.Do(req)
	if err != nil {
		// 超时错误 case : errors.Is(err, context.DeadlineExceeded)
		slog.ErrorContext(ctx, "HTTPClient.Do error", "error", err, "uri", uri)
		return err
	}
	defer resp.Body.Close()

	// 错误状态码返回错误信息
	if err = HTTPResponseCodeError(resp); err != nil {
		slog.ErrorContext(ctx, "HTTPResponseCodeError error", "error", err, "uri", uri)
		return err
	}

	// 创建文件所在目录（如果不存在）
	if err = os.MkdirAll(filepath.Dir(outputpath), 0o664); err != nil {
		slog.ErrorContext(ctx, "os.MkdirAll error", "error", err, "outputpath", outputpath, "uri", uri)
		return err
	}

	outfile, err := os.Create(outputpath)
	if err != nil {
		slog.ErrorContext(ctx, "os.Create error", "error", err, "outputpath", outputpath, "uri", uri)
		return err
	}
	defer outfile.Close()

	_, err = io.Copy(outfile, resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "io.Copy error", "error", err, "outputpath", outputpath, "uri", uri)
		return err
	}

	return nil
}
