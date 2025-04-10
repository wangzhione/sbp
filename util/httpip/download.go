package httpip

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/wangzhione/sbp/chain"
	"github.com/wangzhione/sbp/util/filedir"
)

// Download 下载 uri 到本地文件 outputPath
func Download(ctx context.Context, uri, outputpath string, headerargs ...map[string]string) error {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
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

	resp, err := http.DefaultClient.Do(req)
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

// DownloadIfNotExists 下载文件（如果文件已存在则跳过），失败时清理临时文件
func DownloadIfNotExists(ctx context.Context, uri, outputPath string, headerargs ...map[string]string) (err error) {
	// 如果目标文件已存在，直接跳过
	found, err := filedir.Exist(ctx, outputPath)
	if err != nil {
		return
	}
	if found {
		// 文件存在直接返回
		return
	}

	err = Download(ctx, uri, outputPath, headerargs...)
	if err != nil {
		// 下载失败则尝试删除输出文件
		if rmErr := os.Remove(outputPath); rmErr != nil && !os.IsNotExist(rmErr) {
			slog.WarnContext(ctx, "failed to remove output file", "path", outputPath, "error", rmErr)
		}
		return
	}

	return
}
