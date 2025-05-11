package groupgo

import (
	"log/slog"
	"testing"
	"time"

	"github.com/wangzhione/sbp/chain"
)

// groupgodemo 结构适合你们自行 Ctrl + C -> Ctrl + V 用于实际业务中.
// 需要注意是 context 生命周期, 因为有些 http 服务 call 结束适合, context 会被取消
//
//

var ctx = chain.Context()

func TestDownloadGroup_Download(t *testing.T) {
	chain.InitSLog()

	func() {
		start := time.Now()
		slog.InfoContext(ctx, "Download task start") // "source":"groupgodemo_test.go:18.func1"

		defer func() {
			duration := time.Since(start)
			// 这是个 demo 库, 至少介绍 groupgo 用法
			slog.InfoContext(ctx, "Download task end",
				"duration", duration.Seconds(),
			) // "source":"groupgodemo_test.go:23.1"
		}()
	}()
	/*
	   {"time":"2025-04-15T20:54:41.3282085+08:00","level":"INFO","msg":"Download task start","X-Request-Id":"bf26a46755a84cfeb2ef6dd954f1353d","source":"groupgodemo_test.go:18.func1"}
	   {"time":"2025-04-15T20:54:41.3329293+08:00","level":"INFO","msg":"Download task end","duration":0.0047208,"X-Request-Id":"bf26a46755a84cfeb2ef6dd954f1353d","source":"groupgodemo_test.go:23.1"}
	*/
}
