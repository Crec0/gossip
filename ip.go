package main

import (
	"log"
	"net"
	"net/http"
	"strings"
)

var headersInOrder = []string{
	"x-client-ip",
	"cf-connecting-ip",
	"fastly-client-ip",
	"true-client-ip",
	"x-real-ip",
	"x-cluster-client-ip",
	"x-forwarded",
	"forwarded-for",
	"forwarded",
	"x-appengine-user-ip",
	"cf-pseudo-ipv4",
}

func findIP(r *http.Request) string {
	for _, header := range headersInOrder {
		if ip := r.Header.Get(header); ip != "" {
			log.Printf("Found ip %s in header %s", ip, header)
			return ip
		}
	}

	if ip := r.Header.Get("x-forwarded-for"); ip != "" {
		candidates := strings.Split(ip, ",")
		if len(candidates) > 0 {
			ip = strings.TrimSpace(candidates[0])
		}
		log.Printf("Found ip %s in header X-Forwarded-For", ip)
		return ip
	}

	log.Printf("Fallback to remote header ip %s", r.RemoteAddr)
	if strings.Contains(r.RemoteAddr, ":") {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		return ip
	}
	return r.RemoteAddr
}
