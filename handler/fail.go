package handler

import (
	"fmt"

	"github.com/miekg/dns"
)

var Fail = failHandler{}

type failHandler struct{}

func (failHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	fmt.Println(fmt.Errorf("dns-fuzz: ignored query"))
}
