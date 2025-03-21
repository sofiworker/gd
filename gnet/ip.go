package gnet

import "net"

func IsIPv4(ip string) bool {
	parseIP := net.ParseIP(ip)
	return parseIP != nil && parseIP.To4() != nil
}

func IsIPv6(ip string) bool {
	parseIP := net.ParseIP(ip)
	return parseIP != nil && parseIP.To16() != nil
}
