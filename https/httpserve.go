// Package system provides HTTP server startup, graceful shutdown, and TLS support utilities.
package system

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/wangzhione/sbp/system"
)

// BeginTime 系统启动时间
var BeginTime = time.Now()

// End 主要用于 main 函数中 defer End(context) 操作, 纪录程序结束的行为
func End(ctx context.Context) {
	if cover := recover(); cover != nil {
		// 遇到启动不起来, 异常退出, 打印堆栈方便排除问题
		slog.ErrorContext(ctx, "main init panic error",
			slog.Any("error", cover),
			slog.Time("SystemBeginTime", BeginTime),
			slog.String("GOOS", runtime.GOOS),
			slog.String("BuildVersion", system.BuildVersion),
			slog.String("GitVersion", system.GitVersion),
			slog.String("stack", string(debug.Stack())), // 记录详细的堆栈信息
		)
	}

	end := time.Now()
	slog.InfoContext(ctx, "main init end ...",
		slog.Time("SystemBeginTime", BeginTime),
		slog.Float64("elapsed_hours", end.Sub(BeginTime).Hours()),
		slog.Time("EndTime", end),
		slog.String("GOOS", runtime.GOOS),
		slog.String("BuildVersion", system.BuildVersion),
		slog.String("GitVersion", system.GitVersion),
	)
}

type StopFunc func(ctx context.Context) <-chan struct{}

// ServeLoop 服务启动 loop 主流程
// addr 类似 fmt.Sprintf("0.0.0.0:%d", config.G.Serve.Port) ; 0.0.0.0 默认 ipv4 绑定本机地址
// handler 类似 middleware.MainMiddleware(http.DefaultServeMux)
func ServeLoop(ctx context.Context, addr string, handler http.Handler, stopTime time.Duration, stopmainfunc ...StopFunc) {
	serve := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go ServeShutdown(ctx, serve, stopTime, stopmainfunc...)

	// main server 启动
	slog.InfoContext(ctx, "Server running", slog.String("addr", serve.Addr))
	err := serve.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			slog.InfoContext(ctx, "Server success stop", slog.String("addr", serve.Addr))
			return
		}
		slog.ErrorContext(ctx, "Server ListenAndServe failed error",
			slog.Any("error", err),
			slog.String("addr", serve.Addr),
		)
	}
}

func ServeShutdown(ctx context.Context, server *http.Server, stopTime time.Duration, stopmainfunc ...StopFunc) {
	defer func() {
		if cover := recover(); cover != nil {
			// 遇到启动不起来, 异常退出, 打印堆栈方便排除问题
			slog.ErrorContext(ctx, "Server signal panic error",
				slog.Any("error", cover),
				slog.Time("SystemBeginTime", BeginTime),
				slog.Float64("elapsed_hours", time.Since(BeginTime).Hours()),
				slog.String("stack", string(debug.Stack())), // 记录详细的堆栈信息
			)
		}
	}()

	// 监听系统信号（优雅退出）
	sc := make(chan os.Signal, 1)

	// 监听 Ctrl+C 和 kill or killall 命令
	// 对于 Web 轻量级应用, 花几秒重启代价最小, 还能避免复杂的重度资源处理逻辑
	signal.Notify(sc, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	defer signal.Stop(sc)

	// 等待终止信号
	sig := <-sc
	slog.InfoContext(ctx, "Server Received Shutting down...", "signal", sig)

	// 优雅 stop HTTP 服务器, 设置超时时间的上下文
	timeoutctx, cancel := context.WithTimeout(ctx, stopTime)
	defer cancel()

	// 这部分处理 sig 信号退出
	var stopDone <-chan struct{}
	if len(stopmainfunc) > 0 {
		stopDone = stopmainfunc[0](timeoutctx)
	}

	if err := server.Shutdown(timeoutctx); err != nil {
		slog.ErrorContext(ctx, "Server.Shutdown error", "error", err)
	}
	slog.InfoContext(ctx, "Server gracefully stopped", "SystemBeginTime", BeginTime, "stopTime", stopTime)

	// 等后台任务真正退出（或超时）
	if stopDone != nil {
		select {
		case <-stopDone:
			slog.InfoContext(ctx, "Background tasks stopped")
		case <-timeoutctx.Done():
			slog.WarnContext(ctx, "Background tasks stop timeout", "err", timeoutctx.Err())
		}
	}
}

/*

	1. 方式一

	sudo apt install certbot
	sudo certbot certonly --standalone -d {yourdomain}.com

	/etc/letsencrypt/live/{yourdomain}.com/fullchain.pem  # → certFile
	/etc/letsencrypt/live/{yourdomain}.com/privkey.pem     # → keyFile


	2. 方式二
	openssl req -x509 -newkey rsa:2048 -nodes -keyout server.key -out server.crt -days 365 -subj "/C=CN/ST=Test/L=Dev/O=Local/CN=localhost"
*/

// ServeLoopTLS 服务启动 loop 主流程
// addr 类似 "0.0.0.0:443"
// handler 类似 middleware.MainMiddleware(http.DefaultServeMux)
// 若 certFile 和 keyFile 不为空，则启用 HTTPS
func ServeLoopTLS(ctx context.Context, certFile, keyFile, addr string, handler http.Handler, stopTime time.Duration, stopmainfunc ...StopFunc) {
	serve := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go ServeShutdown(ctx, serve, stopTime, stopmainfunc...)

	slog.InfoContext(ctx, "🔒 HTTPS Server running", slog.String("addr", serve.Addr))
	err := serve.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		if err == http.ErrServerClosed {
			slog.InfoContext(ctx, "HTTPS Server success stop", slog.String("addr", serve.Addr))
			return
		}
		slog.ErrorContext(ctx, "HTTPS Server ListenAndServeTLS failed error",
			slog.Any("error", err),
			slog.String("addr", serve.Addr),
		)
	}
}
