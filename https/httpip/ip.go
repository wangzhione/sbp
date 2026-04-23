package httpip

import (
	"net"
	"net/http"
	"strings"
)

// GetClientIP 获取客户端 ip
func GetClientIP(r *http.Request) (ip string) {
	// X-Forwarded-For: <client>, <proxy1>, <proxy2>
	// X-Forwarded-For (XFF) 在客户端访问服务器的过程中如果需要经过 HTTP 代理或者负载均衡服务器,
	// 可以被用来获取最初发起请求的客户端的 IP 地址, 这个消息首部成为事实上的标准.
	xForwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if len(xForwardedFor) > 0 {
		xForwardedFors := strings.Split(xForwardedFor, ",")
		ip = strings.TrimSpace(xForwardedFors[0])
		if len(ip) > 0 {
			return
		}
	}

	// "X-Real-Ip" nginx 反向代理服务 IP
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if len(ip) > 0 {
		return
	}
	// 兜底直接使用 client 请求的 ip 地址
	ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	return
}
