package httpip

import (
	"io"
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

// GetMACAddress 获取本机 MAC 地址
func GetMACAddress() (macList []string, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, i := range interfaces {
		if i.HardwareAddr != nil {
			macList = append(macList, i.HardwareAddr.String())
		}
	}
	return
}

// GetAllIP 获取本机所有 IP 地址（包括本地和外部）
func GetAllIP(includeIPv6 ...struct{}) ([]string, error) {
	var ipList []string

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if len(includeIPv6) == 0 && ipNet.IP.To4() == nil {
					continue
				}
				ipList = append(ipList, ipNet.IP.String())
			}
		}
	}
	return ipList, nil
}

// GetExternalIP 获取本机外部 IP 地址（排除 127.0.0.1）
func GetExternalIP(includeIPv6 ...struct{}) ([]string, error) {
	var ipList []string

	// 在某些操作系统（如 Linux 或 MacOS）上，获取某些网络接口的地址可能需要管理员权限
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && (len(includeIPv6) == 0 || ipNet.IP.To4() != nil) {
				ipList = append(ipList, ipNet.IP.String())
			}
		}
	}
	return ipList, nil
}

// GetPublicIP 获取本机公网 IP 地址
func GetPublicIP() (string, error) {
	// 公网 IP 查询服务
	resp, err := http.Get("https://api64.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

var privateBlocks = []string{
	"10.",      // 10.0.0.0 – 10.255.255.255
	"172.16.",  // 172.16.0.0 – 172.31.255.255
	"192.168.", // 192.168.0.0 – 192.168.255.255
	"127.",     // Loopback
	"fc00:",    // IPv6 Unique Local Address
	"fe80:",    // IPv6 Link-Local Address
}

// IsPrivateIP 检查 IP 地址是否为内网地址
func IsPrivateIP(ip string) bool {
	for _, block := range privateBlocks {
		if strings.HasPrefix(ip, block) {
			return true
		}
	}
	return false
}
