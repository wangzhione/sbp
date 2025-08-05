package system

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"
)

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
func ServeLoopTLS(ctx context.Context, certFile, keyFile, addr string, handler http.Handler, stopTime time.Duration, stopfunc ...func(context.Context, os.Signal)) {
	serve := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go ServeShutdown(ctx, serve, stopTime, stopfunc...)

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
