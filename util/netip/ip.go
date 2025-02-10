package netip

import (
	"net"
	"net/http"
	"strings"
)

// X-Forwarded-For (XFF) 在客户端访问服务器的过程中如果需要经过 HTTP 代理或者负载均衡服务器,
// 可以被用来获取最初发起请求的客户端的 IP 地址, 这个消息首部成为事实上的标准.
var xForwardedForKeys = [...]string{"X-Forwarded-For", "x-forwarded-for", "X-FORWARDED-FOR"}

// XRealIP nginx 反向代理服务 IP
var xRealIPKeys = [...]string{"X-Real-IP", "X-Real-Ip", "x-real-ip", "X-REAL-IP"}

// GetClientIP 获取客户端 ip
func GetClientIP(r *http.Request) (ip string) {
	// X-Forwarded-For: <client>, <proxy1>, <proxy2>
	for _, xForwardedForKey := range xForwardedForKeys {
		xForwardedFor := strings.TrimSpace(r.Header.Get(xForwardedForKey))
		if len(xForwardedFor) > 0 {
			xForwardedFors := strings.Split(xForwardedFor, ",")
			ip = strings.TrimSpace(xForwardedFors[0])
			if len(ip) > 0 {
				return
			}
			break
		}
	}

	for _, xRealIPKey := range xRealIPKeys {
		ip = strings.TrimSpace(r.Header.Get(xRealIPKey))
		if len(ip) > 0 {
			return
		}
	}

	// 兜底直接使用 client 请求的 ip 地址
	ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	return
}
