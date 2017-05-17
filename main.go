package main

import (
	"fmt"
	"os"

	"github.com/dpb587/dns-fuzz/handler"
	"github.com/miekg/dns"
)

func main() {
	net := "tcp"
	addr := "127.0.0.1:35053"

	if len(os.Args) > 1 {
		net = os.Args[1]

		if len(os.Args) > 2 {
			addr = os.Args[2]
		}
	}

	server := &dns.Server{
		Addr:    addr,
		Net:     net,
		UDPSize: dns.MaxMsgSize,
	}

	dns.Handle("test.", handler.NewDynamicFuzz("test."))
	dns.Handle("custom.example.com.", handler.NewStaticFuzz("size-8.ttl-16.answer.ttl-4.answer"))
	dns.Handle(".", handler.Fail)

	fmt.Println(fmt.Sprintf("listening on %s (%s)", addr, net))

	if err := server.ListenAndServe(); err != nil {
		os.Exit(1)
	}
}
