package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	portToServe, ok := os.LookupEnv("PORT")
	if !ok {
		portToServe = "10059"
	}

	fmt.Printf("Serving on %s\n", portToServe)

	http.HandleFunc("/", serveIP)
	log.Fatal(http.ListenAndServe(":"+portToServe, nil))
}

func serveIP(writer http.ResponseWriter, request *http.Request) {
	ip := net.ParseIP(grabIP(request))
	if ip == nil {
		writer.WriteHeader(400)
		_, _ = fmt.Fprint(writer, "No ip found")
	} else {
		_, _ = fmt.Fprint(writer, ip.To4().String())
	}
}

var ipHeadersInOrder = []string{
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

// Stolen from https://github.com/pbojinov/request-ip
func grabIP(r *http.Request) string {
	fmt.Println(":", r)
	if ip := r.Header.Get("x-client-ip"); ip != "" {
		fmt.Println("triggered: x-client-ip")
		return ip
	}

	if ip := r.Header.Get("x-forwarded-for"); ip != "" {
		candidates := strings.Split(ip, ",")
		if len(candidates) > 0 {
			candidate := strings.TrimSpace(candidates[0])
			if strings.Contains(candidate, ":") {
				candidate, _, _ = net.SplitHostPort(candidate)
			}
			ip = candidate
		}
		fmt.Println("triggered: x-forwarded-for")
		return ip
	}

	for _, headerName := range ipHeadersInOrder {
		if ip := r.Header.Get(headerName); ip != "" {
			fmt.Println("triggered: " + headerName)
			return ip
		}
	}

	fmt.Println("triggered: remote addr")
	return r.RemoteAddr
}
