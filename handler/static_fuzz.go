package handler

import (
	"fmt"

	"github.com/miekg/dns"
)

type StaticFuzz struct {
	fuzz string
}

func NewStaticFuzz(fuzz string) StaticFuzz {
	return StaticFuzz{fuzz: fuzz}
}

func (h StaticFuzz) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	res := &dns.Msg{}
	res.SetReply(req)

	err := Generate(res, h.fuzz)
	if err != nil {
		fmt.Println(err)

		return
	}

	err = w.WriteMsg(res)
	if err != nil {
		fmt.Println(fmt.Errorf("dns-fuzz: failed to write response: %s", err.Error()))
	}
}
