package main

import (
	"context"
	"flag"
	"log"
	"net"
	"strings"

	tcpsrv "lab9_var7/task2_go_tcp_server/server"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:9000", "listen address")
	flag.Parse()

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("tcp server listening on %s", l.Addr().String())

	srv := tcpsrv.NewServer(l, func(_ context.Context, msg string) (string, error) {
		return strings.ToUpper(msg), nil
	})

	if err := srv.Serve(); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
